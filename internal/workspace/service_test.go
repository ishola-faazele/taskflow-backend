package workspace

// import (
// 	"testing"
// )

// func TestWorkspaceService_CreateWorkspace(t *testing.T) {
// 	validOwnerID := "550e8400-e29b-41d4-a716-446655440000" // valid UUID
// 	tests := []struct {
// 		name      string
// 		workspaceName string
// 		ownerID   string
// 		expectError bool
// 	}{
// 		{
// 			name:        "Valid workspace creation",
// 			workspaceName: "Test Workspace",
// 			ownerID:     validOwnerID,
// 			expectError: false,
// 		},
// 		{
// 			name:        "Empty workspace name",
// 			workspaceName: "",
// 			ownerID:     validOwnerID,
// 			expectError: true,
// 		},
// 		{
// 			name:        "Invalid owner ID",
// 			workspaceName: "Test Workspace",
// 			ownerID:     "invalid-uuid",
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			service := &WorkspaceService{
// 				workspaceRepo:   &MockWorkspaceRepository{},
// 				membershipRepo:  &MockMembershipRepository{},
// 				invitationRepo:  &MockInvitationRepository{},
// 			}

// 			workspace, err := service.CreateWorkspace(tt.workspaceName, tt.ownerID)
// 			if tt.expectError && err == nil {
// 				t.Errorf("expected error but got none")
// 			}
// 			if !tt.expectError && err != nil {
// 				t.Errorf("did not expect error but got: %v", err)
// 			}
// 			if !tt.expectError && workspace.Name != tt.workspaceName {
// 				t.Errorf("expected workspace name %s, got %s", tt.workspaceName, workspace.Name)
// 			}
// 		})
// 	}

// }