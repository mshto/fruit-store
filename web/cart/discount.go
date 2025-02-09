package cart

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/mshto/fruit-store/cache"
	"github.com/mshto/fruit-store/entity"
	"github.com/mshto/fruit-store/repository"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// AddDiscout add user discout
func (ph cartHandler) AddDiscout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		ph.log.Errorf("failed to get user uuid, user: %v, error: %v", ctx.Value(middleware.UserUUID).(string), err)
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	dsc := &entity.Discount{}
	err = json.NewDecoder(r.Body).Decode(dsc)
	if err != nil {
		ph.log.Errorf("failed to decode discount, user: %v, error: %v", userUUID, err)
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	sale, err := ph.bil.GetDiscountByUser(userUUID)
	if err != nil && err != cache.ErrNotFound {
		ph.log.Errorf("failed to get discount by user, user: %v, error: %v", userUUID, err)
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}
	if sale.ID != "" {
		ph.log.Warnf("discount is already added, user: %v", userUUID)
		response.RenderFailedResponse(w, http.StatusConflict, errors.New("discount is already added"))
		return
	}

	dscRepo, err := ph.discRepo.GetDiscount(dsc.ID)
	if err == repository.ErrNotFound {
		ph.log.Errorf("failed to get discount, user: %v, error: %v", userUUID, err)
		response.RenderFailedResponse(w, http.StatusNotFound, err)
		return
	}
	if err != nil {
		ph.log.Errorf("failed to get discount, user: %v, error: %v", userUUID, err)
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = ph.bil.SetDiscount(userUUID, dscRepo)
	if err != nil {
		ph.log.Errorf("failed to set discount, user: %v, error: %v", userUUID, err)
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}
