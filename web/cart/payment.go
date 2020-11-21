package cart

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/mshto/fruit-store/entity"
	"github.com/mshto/fruit-store/web/common/response"
	"github.com/mshto/fruit-store/web/middleware"
)

// GetAllProducts retrieves all products from db
func (ph cartHandler) AddPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userUUID, err := uuid.Parse(ctx.Value(middleware.UserUUID).(string))
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	pmt := entity.Payment{}
	err = json.NewDecoder(r.Body).Decode(&pmt)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = ph.bil.ValidateCard(pmt)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}

	err = ph.cartRepo.RemoveUserProducts(userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	err = ph.bil.RemoveDiscount(userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}

	response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}
