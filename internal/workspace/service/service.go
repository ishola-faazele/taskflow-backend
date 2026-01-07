package workspace

import (
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/ishola-faazele/taskflow/internal/utils/jwt"
	. "github.com/ishola-faazele/taskflow/internal/workspace/entity"
	. "github.com/ishola-faazele/taskflow/internal/workspace/repository"
	"github.com/ishola-faazele/taskflow/pkg/utils"
	"github.com/ishola-faazele/taskflow/pkg/utils/domain_errors"
)

type WorkspaceService struct {
	WorkspaceRepo  WorkspaceRepository
	MembershipRepo MembershipRepository
	InvitationRepo InvitationRepository
	jwtUtil        *jwt.JWTUtils
	emailService   *utils.EmailService
}

func NewWorkspaceService(workspaceRepo WorkspaceRepository, invitationRepo InvitationRepository, membershipRepo MembershipRepository) *WorkspaceService {
	emailConfig := utils.EmailConfig{
		SMTPHost:    os.Getenv("SMTP_HOST"),
		SMTPPort:    "587",
		SenderEmail: os.Getenv("SMTP_USER"),
		SenderName:  "TaskFlow Support",
		AppPassword: os.Getenv("SMTP_PASS"),
		FrontendURL: "http://localhost:3000",
	}
	return &WorkspaceService{
		WorkspaceRepo:  workspaceRepo,
		MembershipRepo: membershipRepo,
		InvitationRepo: invitationRepo,
		jwtUtil:        jwt.NewJWTUtils(jwt.DefaultTokenConfig()),
		emailService:   utils.NewEmailService(emailConfig),
	}
}

func (s *WorkspaceService) CreateWorkspace(name, ownerID string) (*Workspace, domain_errors.DomainError) {
	// check for valid name and ownerID
	if name == "" {
		return nil, domain_errors.NewValidationErrorWithValue("name", name, "EMPTY WORKSPACE NAME")
	}

	workspace := &Workspace{
		ID:        uuid.NewString(),
		Name:      name,
		OwnerID:   ownerID,
		CreatedAt: time.Now().UTC(),
	}
	return s.WorkspaceRepo.Create(workspace)
}

func (s *WorkspaceService) GetWorkspaceByID(id string) (*Workspace, domain_errors.DomainError) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("id", id, "WORKSPACE_ID IS NOT A VALID UUID")
	}
	return s.WorkspaceRepo.GetByID(id)
}

func (s *WorkspaceService) UpdateWorkspace(id, name, requester string) (*Workspace, domain_errors.DomainError) {
	if name == "" {
		return nil, domain_errors.NewValidationErrorWithValue("name", name, "EMPTY WORKSPACE NAME")
	}
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspaceID", id, "WORKSPACE ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(requester); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("requester_id", requester, "REQUESTER_ID IS NOT A VALID UUID")
	}

	workspace, err := s.GetWorkspaceByID(id)
	if err != nil {
		return nil, err
	}
	if requester != workspace.OwnerID {
		return nil, domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF WORKSPACE")
	}
	workspace.Name = name

	return s.WorkspaceRepo.Update(workspace)
}

func (s *WorkspaceService) DeleteWorkspace(id, requester string) domain_errors.DomainError {
	if err := uuid.Validate(id); err != nil {
		return domain_errors.NewValidationErrorWithValue("workspaceID", id, "WorkspaceID is not a valid UUID")
	}
	ws, err := s.GetWorkspaceByID(id)
	if err != nil {
		return err
	}
	if ws.OwnerID != requester {
		return domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF WORKSPACE")
	}

	return s.WorkspaceRepo.Delete(id)
}

func (s *WorkspaceService) ListWorkspacesByOwner(ownerID string) ([]*Workspace, domain_errors.DomainError) {
	log.Println("USER ID: ", ownerID)
	if err := uuid.Validate(ownerID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("ownerID", ownerID, "OWNER ID IS NOT A VALID UUID")
	}
	return s.WorkspaceRepo.ListByOwner(ownerID)

}

// INVITATION FUNCTIONS
func (s *WorkspaceService) CreateInvitation(invitee, inviter, ws, email string, role Role) (*Invitation, domain_errors.DomainError) {
	// validate inputs
	if err := uuid.Validate(invitee); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("invitee_id", invitee, "INVITEE_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(inviter); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("inviter_id", invitee, "INVITER_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(ws); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspace_id", ws, "WorkspaceID is not a valid UUID")
	}
	if !utils.IsValidEmail(email) {
		return nil, domain_errors.NewValidationErrorWithValue("email", email, "INVALID EMAIL FORMAT")
	}

	inv := &Invitation{
		ID:           uuid.NewString(),
		InviteeID:    invitee,
		InviterID:    inviter,
		InviteeEmail: email,
		WorkspaceID:  ws,
		Role:         role,
		IsValid:      true,
		CreatedAt:    time.Now().UTC(),
	}
	// send invitation link to email
	inv, err := s.InvitationRepo.Create(inv)
	if err != nil {
		return nil, err
	}
	token, errToken := s.jwtUtil.GenerateInvitationToken(inv.ID, ws, inv.InviterID, inv.InviteeEmail, inv.InviteeID, string(role))
	if errToken != nil {
		return nil, domain_errors.NewInternalError("FAILED GENERATING INVITATION TOKEN", errToken)
	}
	errEmail := s.emailService.SendInvitationLink(inv.InviteeEmail, inv.WorkspaceID, string(role), token, "/api/workspace/invitation/verify?token")
	if errEmail != nil {
		return nil, domain_errors.NewInternalError("FAILED SENDING INVITATION EMAIL", errEmail)
	}
	return inv, nil
}

func (s *WorkspaceService) GetInvitation(id string) (*Invitation, error) {
	if err := uuid.Validate(id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("invitation_id", id, "INVITATION_ID IS NOT A VALID UUID")
	}
	return s.InvitationRepo.GetByID(id)
}

func (s *WorkspaceService) DeleteInvitation(id, requester string) error {
	if err := uuid.Validate(id); err != nil {
		return domain_errors.NewValidationErrorWithValue("invitation_id", id, "INVITATION_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(requester); err != nil {
		return domain_errors.NewValidationErrorWithValue("requester_id", requester, "REQUESTER_ID IS NOT A VALID UUID")
	}
	// Check requester owns invitation
	invitation, err := s.GetInvitation(id)
	if err != nil {
		return err
	}
	if invitation.InviterID != requester {
		return domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE INVITER OF INVITATION")
	}
	return s.InvitationRepo.DeleteInvitation(id)
}
func (s *WorkspaceService) ListWorkspaceInvitations(ws_id, requester string) ([]*Invitation, domain_errors.DomainError) {
	if err := uuid.Validate(ws_id); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspace_id", ws_id, "WORKSPACE_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(requester); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("requester_id", requester, "REQUESTER_ID IS NOT A VALID UUID")
	}
	// Check requester owns workspace
	ws, err := s.GetWorkspaceByID(ws_id)
	if err != nil {
		return nil, err
	}
	if ws.OwnerID != requester {
		return nil, domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF WORKSPACE")
	}
	return s.InvitationRepo.ListInvitationToWorkspace(ws_id)
}

// MEMBERSHIP FUNCTIONS

func (s *WorkspaceService) AddMembership(token string) (*Membership, error) {
	claims, err := s.jwtUtil.ParseInvitationToken(token)
	if err != nil {
		return nil, domain_errors.NewUnauthorizedError("INVALID OR EXPIRED INVITATION TOKEN")
	}
	userID := claims.InviteeID
	ws_id := claims.WorkspaceID
	role := Role(claims.Role)
	membership := &Membership{
		UserID:      userID,
		WorkspaceID: ws_id,
		Role:        role,
		CreatedAt:   time.Now().UTC(),
	}
	return s.MembershipRepo.Add(membership)
}
func (s *WorkspaceService) RemoveMembership(userID, workspaceID, requester string) error {
	// validate inputs
	if err := uuid.Validate(userID); err != nil {
		return domain_errors.NewValidationErrorWithValue("user_id", userID, "USER_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(workspaceID); err != nil {
		return domain_errors.NewValidationErrorWithValue("workspace_id", workspaceID, "WORKSPACE_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(requester); err != nil {
		return domain_errors.NewValidationErrorWithValue("requester_id", requester, "REQUESTER_ID IS NOT A VALID UUID")
	}
	// check requester owns workspace
	ws, err := s.GetWorkspaceByID(workspaceID)
	if err != nil {
		return err
	}
	if ws.OwnerID != requester {
		return domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF WORKSPACE")
	}
	return s.MembershipRepo.Remove(userID, workspaceID)
}

func (s *WorkspaceService) ListWorkspaceMembers(workspaceID, requester string) ([]*Membership, error) {
	// validate inputs
	if err := uuid.Validate(workspaceID); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("workspace_id", workspaceID, "WORKSPACE_ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(requester); err != nil {
		return nil, domain_errors.NewValidationErrorWithValue("requester_id", requester, "REQUESTER_ID IS NOT A VALID UUID")
	}
	// check requester owns workspace
	ws, err := s.GetWorkspaceByID(workspaceID)
	if err != nil {
		return nil, err
	}
	if ws.OwnerID != requester {
		return nil, domain_errors.NewUnauthorizedError("REQUESTER IS NOT THE OWNER OF WORKSPACE")
	}
	return s.MembershipRepo.ListByWorkspace(workspaceID)
}
func (s *WorkspaceService) IsMember(userID, workspaceID string) (bool, error) {
	if err := uuid.Validate(userID); err != nil {
		return false, domain_errors.NewValidationErrorWithValue("user_id", userID, "USER ID IS NOT A VALID UUID")
	}
	if err := uuid.Validate(workspaceID); err != nil {
		return false, domain_errors.NewValidationErrorWithValue("user_id", userID, "USER ID IS NOT A VALID UUID")
	}
	return s.MembershipRepo.IsMember(userID, workspaceID)
}
