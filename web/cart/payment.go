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

	err = ph.repo.Cart.RemoveUserProducts(userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusInternalServerError, err)
		return
	}
	err = ph.bil.Pay(userUUID)
	if err != nil {
		response.RenderFailedResponse(w, http.StatusBadRequest, err)
		return
	}
	// if err != nil && err != cache.ErrNotFound {
	// 	response.RenderFailedResponse(w, http.StatusInternalServerError, err)
	// 	return
	// }
	// if sale.ID != "" {
	// 	response.RenderFailedResponse(w, http.StatusConflict, errors.New("discount is already added"))
	// 	return
	// }

	// dscRepo, err := ph.repo.Discount.GetDiscount(dsc.ID)
	// if err == repository.ErrNotFound {
	// 	response.RenderFailedResponse(w, http.StatusNotFound, err)
	// 	return
	// }
	// if err != nil {
	// 	response.RenderFailedResponse(w, http.StatusInternalServerError, err)
	// 	return
	// }

	// err = ph.bil.SetDiscount(userUUID, dscRepo)
	// if err != nil {
	// 	response.RenderFailedResponse(w, http.StatusInternalServerError, err)
	// 	return
	// }

	// response.RenderResponse(w, http.StatusCreated, response.EmptyResp{})
}
