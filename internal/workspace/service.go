package workspace


type WorkspaceService struct {
	workspaceRepo   WorkspaceRepository
	membershipRepo  MembershipRepository
	invitationRepo  InvitationRepository
}

func (s *WorkspaceService) CreateWorkspace(name, ownerID string) (*Workspace, error) {
	workspace := &Workspace{
		Name:    name,
		OwnerID: ownerID,
	}
	return s.workspaceRepo.Create(workspace)
}