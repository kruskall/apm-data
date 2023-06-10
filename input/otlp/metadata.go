// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package otlp

import (
	"fmt"
	"net/netip"
	"regexp"
	"strconv"
	"strings"

	"github.com/elastic/apm-data/model/modelpb"
	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/collector/semconv/v1.5.0"
)

const (
	AgentNameJaeger = "Jaeger"
)

var (
	serviceNameInvalidRegexp = regexp.MustCompile("[^a-zA-Z0-9 _-]")
)

func translateResourceMetadata(resource pcommon.Resource, out *modelpb.APMEvent) {
	var exporterVersion string
	resource.Attributes().Range(func(k string, v pcommon.Value) bool {
		switch k {
		// service.*
		case semconv.AttributeServiceName:
			out.Service.Name = cleanServiceName(v.Str())
		case semconv.AttributeServiceVersion:
			out.Service.Version = truncate(v.Str())
		case semconv.AttributeServiceInstanceID:
			out.Service.Node.Name = truncate(v.Str())

		// deployment.*
		case semconv.AttributeDeploymentEnvironment:
			out.Service.Environment = truncate(v.Str())

		// telemetry.sdk.*
		case semconv.AttributeTelemetrySDKName:
			out.Agent.Name = truncate(v.Str())
		case semconv.AttributeTelemetrySDKVersion:
			out.Agent.Version = truncate(v.Str())
		case semconv.AttributeTelemetrySDKLanguage:
			out.Service.Language.Name = truncate(v.Str())

		// cloud.*
		case semconv.AttributeCloudProvider:
			out.Cloud.Provider = truncate(v.Str())
		case semconv.AttributeCloudAccountID:
			out.Cloud.AccountId = truncate(v.Str())
		case semconv.AttributeCloudRegion:
			out.Cloud.Region = truncate(v.Str())
		case semconv.AttributeCloudAvailabilityZone:
			out.Cloud.AvailabilityZone = truncate(v.Str())
		case semconv.AttributeCloudPlatform:
			out.Cloud.ServiceName = truncate(v.Str())

		// container.*
		case semconv.AttributeContainerName:
			out.Container.Name = truncate(v.Str())
		case semconv.AttributeContainerID:
			out.Container.Id = truncate(v.Str())
		case semconv.AttributeContainerImageName:
			out.Container.ImageName = truncate(v.Str())
		case semconv.AttributeContainerImageTag:
			out.Container.ImageTag = truncate(v.Str())
		case "container.runtime":
			out.Container.Runtime = truncate(v.Str())

		// k8s.*
		case semconv.AttributeK8SNamespaceName:
			out.Kubernetes.Namespace = truncate(v.Str())
		case semconv.AttributeK8SNodeName:
			out.Kubernetes.NodeName = truncate(v.Str())
		case semconv.AttributeK8SPodName:
			out.Kubernetes.PodName = truncate(v.Str())
		case semconv.AttributeK8SPodUID:
			out.Kubernetes.PodUid = truncate(v.Str())

		// host.*
		case semconv.AttributeHostName:
			out.Host.Hostname = truncate(v.Str())
		case semconv.AttributeHostID:
			out.Host.Id = truncate(v.Str())
		case semconv.AttributeHostType:
			out.Host.Type = truncate(v.Str())
		case "host.arch":
			out.Host.Architecture = truncate(v.Str())

		// process.*
		case semconv.AttributeProcessPID:
			out.Process.Pid = uint32(v.Int())
		case semconv.AttributeProcessCommandLine:
			out.Process.CommandLine = truncate(v.Str())
		case semconv.AttributeProcessExecutablePath:
			out.Process.Executable = truncate(v.Str())
		case "process.runtime.name":
			out.Service.Runtime.Name = truncate(v.Str())
		case "process.runtime.version":
			out.Service.Runtime.Version = truncate(v.Str())

		// os.*
		case semconv.AttributeOSType:
			out.Host.Os.Platform = strings.ToLower(truncate(v.Str()))
		case semconv.AttributeOSDescription:
			out.Host.Os.Full = truncate(v.Str())
		case semconv.AttributeOSName:
			out.Host.Os.Name = truncate(v.Str())
		case semconv.AttributeOSVersion:
			out.Host.Os.Version = truncate(v.Str())

		// device.*
		case semconv.AttributeDeviceID:
			out.Device.Id = truncate(v.Str())
		case semconv.AttributeDeviceModelIdentifier:
			out.Device.Model.Identifier = truncate(v.Str())
		case semconv.AttributeDeviceModelName:
			out.Device.Model.Name = truncate(v.Str())
		case "device.manufacturer":
			out.Device.Manufacturer = truncate(v.Str())

		// Legacy OpenCensus attributes.
		case "opencensus.exporterversion":
			exporterVersion = v.Str()

		// timestamp attribute to deal with time skew on mobile
		// devices. APM server should drop this field.
		case "telemetry.sdk.elastic_export_timestamp":
			// Do nothing.

		default:
			if out.Labels == nil {
				out.Labels = make(modelpb.Labels)
			}
			if out.NumericLabels == nil {
				out.NumericLabels = make(modelpb.NumericLabels)
			}
			setLabel(replaceDots(k), out, ifaceAttributeValue(v))
		}
		return true
	})

	// https://www.elastic.co/guide/en/ecs/current/ecs-os.html#field-os-type:
	//
	// "One of these following values should be used (lowercase): linux, macos, unix, windows.
	// If the OS you’re dealing with is not in the list, the field should not be populated."
	switch out.GetHost().GetOs().GetPlatform() {
	case "windows", "linux":
		out.Host.Os.Type = out.Host.Os.Platform
	case "darwin":
		out.Host.Os.Type = "macos"
	case "aix", "hpux", "solaris":
		out.Host.Os.Type = "unix"
	}

	switch out.GetHost().GetOs().GetName() {
	case "Android":
		out.Host.Os.Type = "android"
	case "iOS":
		out.Host.Os.Type = "ios"
	}

	if strings.HasPrefix(exporterVersion, "Jaeger") {
		// version is of format `Jaeger-<agentlanguage>-<version>`, e.g. `Jaeger-Go-2.20.0`
		const nVersionParts = 3
		versionParts := strings.SplitN(exporterVersion, "-", nVersionParts)
		if out.Service.Language.Name == "" && len(versionParts) == nVersionParts {
			out.Service.Language.Name = versionParts[1]
		}
		if v := versionParts[len(versionParts)-1]; v != "" {
			out.Agent.Version = v
		}
		out.Agent.Name = AgentNameJaeger

		// Translate known Jaeger labels.
		if clientUUID, ok := out.Labels["client-uuid"]; ok {
			out.Agent.EphemeralId = clientUUID.Value
			delete(out.Labels, "client-uuid")
		}
		if systemIP, ok := out.Labels["ip"]; ok {
			if ip, err := netip.ParseAddr(systemIP.Value); err == nil {
				out.Host.Ip = []string{ip.String()}
			}
			delete(out.Labels, "ip")
		}
	}

	if out.GetService().GetName() == "" {
		// service.name is a required field.
		if out.Service == nil {
			out.Service = &modelpb.Service{}
		}
		out.Service.Name = "unknown"
	}
	if out.GetAgent().GetName() == "" {
		// agent.name is a required field.
		if out.Agent == nil {
			out.Agent = &modelpb.Agent{}
		}
		out.Agent.Name = "otlp"
	}
	if out.GetAgent().GetVersion() == "" {
		// agent.version is a required field.
		if out.Agent == nil {
			out.Agent = &modelpb.Agent{}
		}
		out.Agent.Version = "unknown"
	}
	if out.GetService().GetLanguage().GetName() != "" {
		out.Agent.Name = fmt.Sprintf("%s/%s", out.Agent.Name, out.Service.Language.Name)
	} else {
		if out.Service == nil {
			out.Service = &modelpb.Service{}
		}
		if out.Service.Language == nil {
			out.Service.Language = &modelpb.Language{}
		}
		out.Service.Language.Name = "unknown"
	}

	// Set the decoded labels as "global" -- defined at the service level.
	for k, v := range out.Labels {
		v.Global = true
		out.Labels[k] = v
	}
	for k, v := range out.NumericLabels {
		v.Global = true
		out.NumericLabels[k] = v
	}
}

func cleanServiceName(name string) string {
	return serviceNameInvalidRegexp.ReplaceAllString(truncate(name), "_")
}

func ifaceAttributeValue(v pcommon.Value) interface{} {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		return truncate(v.Str())
	case pcommon.ValueTypeBool:
		return strconv.FormatBool(v.Bool())
	case pcommon.ValueTypeInt:
		return float64(v.Int())
	case pcommon.ValueTypeDouble:
		return v.Double()
	case pcommon.ValueTypeSlice:
		return ifaceAttributeValueSlice(v.Slice())
	}
	return nil
}

func ifaceAttributeValueSlice(slice pcommon.Slice) []interface{} {
	values := make([]interface{}, slice.Len())
	for i := range values {
		values[i] = ifaceAttributeValue(slice.At(i))
	}
	return values
}

// initEventLabels initializes an event-specific labels from an event.
func initEventLabels(e *modelpb.APMEvent) {
	e.Labels = modelpb.Labels(e.Labels).Clone()
	e.NumericLabels = modelpb.NumericLabels(e.NumericLabels).Clone()
}

func setLabel(key string, event *modelpb.APMEvent, v interface{}) {
	switch v := v.(type) {
	case string:
		modelpb.Labels(event.Labels).Set(key, v)
	case bool:
		modelpb.Labels(event.Labels).Set(key, strconv.FormatBool(v))
	case float64:
		modelpb.NumericLabels(event.NumericLabels).Set(key, v)
	case int64:
		modelpb.NumericLabels(event.NumericLabels).Set(key, float64(v))
	case []interface{}:
		if len(v) == 0 {
			return
		}
		switch v[0].(type) {
		case string:
			value := make([]string, len(v))
			for i := range v {
				value[i] = v[i].(string)
			}
			modelpb.Labels(event.Labels).SetSlice(key, value)
		case float64:
			value := make([]float64, len(v))
			for i := range v {
				value[i] = v[i].(float64)
			}
			modelpb.NumericLabels(event.NumericLabels).SetSlice(key, value)
		}
	}
}
