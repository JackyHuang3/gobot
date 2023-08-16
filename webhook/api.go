package webhook

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"gobot/global"

	"github.com/julienschmidt/httprouter"
)

var (
	vStore     global.IStore
	vConfig    global.IConfig
	LogPrintln func(datas ...interface{})
	LogPrintf  func(format string, datas ...interface{})
)

func Init(config global.IConfig, store global.IStore) {
	vConfig = config
	vStore = store
}

func RouterInit(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/webhook", webHookVerify)
	router.HandlerFunc(http.MethodPost, "/webhook", webHookHandle)
}

func webHookVerify(w http.ResponseWriter, r *http.Request) {
	mode := WebParams(r).Get("hub.mode")
	token := WebParams(r).Get("hub.verify_token")
	challenge := WebParams(r).Get("hub.challenge")

	if mode == "subscribe" && token == vConfig.Get(global.C_VerifyToken) {
		WebResponse(w, r, []byte(challenge), http.StatusOK)
	} else {
		WebResponse(w, r, []byte{}, http.StatusForbidden)
	}
}

func webHookHandle(w http.ResponseWriter, r *http.Request) {
	reqBuf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WebResponse(w, r, []byte(err.Error()), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	LogPrintf("webHookHandle receive message: %s", string(reqBuf))

	// prepare for entry
	wd := struct {
		Object string
		Entry  []json.RawMessage
	}{}
	if err := json.Unmarshal(reqBuf, &wd); err != nil {
		WebResponse(w, r, []byte(err.Error()), http.StatusInternalServerError)
		return
	}
	var handleFunc func(buf json.RawMessage) error
	switch wd.Object {
	case "page":
		handleFunc = handlePageWebHook
	case "payments":
		handleFunc = handlePaymentWebHook
	default:
		WebResponse(w, r, []byte("unsupport object: "+wd.Object), http.StatusInternalServerError)
		return
	}

	// handle
	for _, entry := range wd.Entry {
		if handleFunc != nil {
			if err := handleFunc(entry); err != nil {
				WebResponse(w, r, []byte(err.Error()), http.StatusInternalServerError)
				return
			}
		}
	}

	// response
	WebResponse(w, r, []byte("successed"), http.StatusOK)
}

func handlePageWebHook(buf json.RawMessage) error {
	entry := struct {
		ID        string
		Time      uint64
		Messaging []json.RawMessage
	}{}
	if err := json.Unmarshal(buf, &entry); err != nil {
		return err
	}
	if entry.ID != vConfig.Get(global.C_PageID) {
		return errors.New("permission denied")
	}

	// handle message
	for _, msgBuf := range entry.Messaging {
		if handle, err := MsgHandle(msgBuf); err != nil || handle {
			if handle {
				break
			}
			return err
		}
		if handle, err := FeedbackHandle(msgBuf); err != nil || handle {
			if handle {
				break
			}
			return err
		}
	}
	return nil
}

func handlePaymentWebHook(buf json.RawMessage) error {
	entry := struct {
		ID            string
		Time          uint64
		ChangedFields []string
	}{}
	if err := json.Unmarshal(buf, &entry); err != nil {
		return err
	}

	// handle message
	for _, actions := range entry.ChangedFields {
		if handle, err := PaymentHandle(entry.ID, actions); err != nil || handle {
			if handle {
				break
			}
			return err
		}
	}
	return nil
}
