// Package v1 provides a bridge to the actual generated protobuf types
package v1

// Re-export all types from the actual generated location
import (
	clawv1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/claw/v1"
)

// Source types
type SourceData = clawv1.SourceData
type SourceSchedule = clawv1.SourceSchedule

// Request types
type CreateSourceRequest = clawv1.CreateSourceRequest
type GetSourceRequest = clawv1.GetSourceRequest
type UpdateSourceRequest = clawv1.UpdateSourceRequest
type DeleteSourceRequest = clawv1.DeleteSourceRequest
type ListSourcesRequest = clawv1.ListSourcesRequest

// Response types
type CreateSourceResponse = clawv1.CreateSourceResponse
type GetSourceResponse = clawv1.GetSourceResponse
type UpdateSourceResponse = clawv1.UpdateSourceResponse
type DeleteSourceResponse = clawv1.DeleteSourceResponse
type ListSourcesResponse = clawv1.ListSourcesResponse

// Helper types
type SourceScheduleList = clawv1.SourceScheduleList

// Re-export the file descriptor
var File_claw_v1_source_proto = clawv1.File_claw_v1_source_proto
var File_claw_v1_source_service_proto = clawv1.File_claw_v1_source_service_proto