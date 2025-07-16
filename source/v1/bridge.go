// Package v1 provides a bridge to the actual generated protobuf types
package v1

// Re-export all types from the actual generated location
import (
	sourcev1 "github.com/tigorlazuardi/claw/lib/claw/gen/proto/source/v1"
)

// Source types
type Source = sourcev1.Source
type Schedule = sourcev1.Schedule

// Request types
type CreateSourceRequest = sourcev1.CreateSourceRequest
type GetSourceRequest = sourcev1.GetSourceRequest
type UpdateSourceRequest = sourcev1.UpdateSourceRequest
type DeleteSourceRequest = sourcev1.DeleteSourceRequest
type ListSourcesRequest = sourcev1.ListSourcesRequest

// Response types
type CreateSourceResponse = sourcev1.CreateSourceResponse
type GetSourceResponse = sourcev1.GetSourceResponse
type UpdateSourceResponse = sourcev1.UpdateSourceResponse
type DeleteSourceResponse = sourcev1.DeleteSourceResponse
type ListSourcesResponse = sourcev1.ListSourcesResponse

// Helper types
type SourceSchedules = sourcev1.SourceSchedules

// Re-export the file descriptor
var File_source_v1_source_proto = sourcev1.File_source_v1_source_proto
var File_source_v1_source_service_proto = sourcev1.File_source_v1_source_service_proto