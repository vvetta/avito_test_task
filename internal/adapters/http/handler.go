package http

import (
	"errors"
	"net/http"
	"encoding/json"

	"avito_test_task/internal/usecase"
	"avito_test_task/internal/adapters/http/openapi"
	"avito_test_task/internal/domain"
)

type Handler struct {
	teams usecase.TeamService
	users usecase.UserService
	prs usecase.PullRequestService
}

func NewHandler(users usecase.UserService, teams usecase.TeamService, prs usecase.PullRequestService) *Handler {
	return &Handler{
		teams: teams,
		users: users,
		prs: prs,
	}
}

func (h *Handler) PostTeamAdd(w http.ResponseWriter, r *http.Request) {
var req openapi.Team
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	team := fromAPITeam(req)

	result, err := h.teams.CreateTeam(r.Context(), team)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTeamExists):
			writeError(w, http.StatusBadRequest, string(openapi.TEAMEXISTS), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	resp := toAPITeam(result)
	writeJSON(w, http.StatusCreated, map[string]any{
		"team": resp,
	})
}

func (h *Handler) GetTeamGet(w http.ResponseWriter, r *http.Request, params openapi.GetTeamGetParams) {
	team, err := h.teams.GetTeam(r.Context(), params.TeamName)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, string(openapi.NOTFOUND), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	resp := toAPITeam(team)
	writeJSON(w, http.StatusOK, resp)
}

func (h *Handler) PostUsersSetIsActive(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostUsersSetIsActiveJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	user, err := h.users.SetUserActive(r.Context(), req.UserId, req.IsActive)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUserNotFound):
			writeError(w, http.StatusNotFound, string(openapi.NOTFOUND), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	respUser := toAPIUser(user)
	writeJSON(w, http.StatusOK, map[string]any{
		"user": respUser,
	})
}

func (h *Handler) GetUsersGetReview(w http.ResponseWriter, r *http.Request, params openapi.GetUsersGetReviewParams) {
	prs, err := h.users.GetUserReviews(r.Context(), params.UserId)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		return
	}

	short := toAPIPullRequestShorts(prs)

	writeJSON(w, http.StatusOK, map[string]any{
		"user_id":       params.UserId,
		"pull_requests": short,
	})
}

func (h *Handler) PostPullRequestCreate(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostPullRequestCreateJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	pr := domain.PullRequest{
		PullRequestID:   req.PullRequestId,
		PullRequestName: req.PullRequestName,
		AuthorID:        req.AuthorId,
	}

	created, err := h.prs.CreatePR(r.Context(), pr)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPullRequestExists):
			writeError(w, http.StatusConflict, string(openapi.PREXISTS), err.Error())
		case errors.Is(err, domain.ErrUserNotFound),
			errors.Is(err, domain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, string(openapi.NOTFOUND), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	resp := toAPIPullRequest(created)
	writeJSON(w, http.StatusCreated, map[string]any{
		"pr": resp,
	})
}

func (h *Handler) PostPullRequestMerge(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostPullRequestMergeJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	pr, err := h.prs.MergePR(r.Context(), req.PullRequestId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPullRequestNotFound):
			writeError(w, http.StatusNotFound, string(openapi.NOTFOUND), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	resp := toAPIPullRequest(pr)
	writeJSON(w, http.StatusOK, map[string]any{
		"pr": resp,
	})
}

func (h *Handler) PostPullRequestReassign(w http.ResponseWriter, r *http.Request) {
	var req openapi.PostPullRequestReassignJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BODY", "invalid request body")
		return
	}

	pr, replacedBy, err := h.prs.ReassignReviewer(r.Context(), req.PullRequestId, req.OldUserId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrPullRequestNotFound),
			errors.Is(err, domain.ErrUserNotFound),
			errors.Is(err, domain.ErrTeamNotFound):
			writeError(w, http.StatusNotFound, string(openapi.NOTFOUND), err.Error())
		case errors.Is(err, domain.ErrPullRequestMerged):
			writeError(w, http.StatusConflict, string(openapi.PRMERGED), err.Error())
		case errors.Is(err, domain.ErrReviewerNotAssigned):
			writeError(w, http.StatusConflict, string(openapi.NOTASSIGNED), err.Error())
		case errors.Is(err, domain.ErrNoCandidate):
			writeError(w, http.StatusConflict, string(openapi.NOCANDIDATE), err.Error())
		default:
			writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
		}
		return
	}

	respPR := toAPIPullRequest(pr)
	writeJSON(w, http.StatusOK, map[string]any{
		"pr":          respPR,
		"replaced_by": replacedBy,
	})
}
