package controllers

import (
	"github.com/revel/revel"
)

type LambdaReq struct {
	ID     string `json:"device_id"`
	Intent string `json:"intent"`
	Data   string `json:"data"`
}
type LambdaResp struct {
	Data string `json:"data"`
}

type TohelpDesk struct {
	Command string `json:"command"`
	Room    int    `json:"room"`
	Data    string `json:"data"`
}

type App struct {
	*revel.Controller
}

var deskChannel = make(chan LambdaReq, 1)
var lambdaRespChan = make(chan bool, 1)

func (c App) Index() revel.Result {
	return c.Render()
}
func (c App) LamdaApi() revel.Result {
	lambdaReq := LambdaReq{}
	resp := LambdaResp{}
	if err := c.Params.BindJSON(&lambdaReq); err != nil {
		c.Response.Status = 400
		revel.ERROR.Println("Invalid request")
		resp.Data = "Invalid request"
		return c.RenderJSON(resp)
	}
	revel.INFO.Println("Lambda request received: ", &lambdaReq)
	switch lambdaReq.Intent {
	case "service":
		deskChannel <- lambdaReq
		if sent := <-lambdaRespChan; sent {
			resp.Data = `The fornt desk has recived your request successfully,
			 You can expect delivery of " + lambdaReq.Data + "shortly`
		}
	case "checkout":
		deskChannel <- lambdaReq
		if sent := <-lambdaRespChan; sent {
			resp.Data = `The fornt desk has recived your request for Checkout,
			 You can visit the reception and return the keys of your room `
		}
		resp.Data = ""
	case "registerEcho":
		deskChannel <- lambdaReq
		if sent := <-lambdaRespChan; sent {
			resp.Data = "ECHO device successfully registerd to room number " + lambdaReq.Data
		}

	default:
		resp.Data = ""
	}
	return c.RenderJSON(resp)
}

func (c App) AlexaOrder(data string, ws revel.ServerWebSocket) revel.Result {
	revel.INFO.Println(" Web socket Connection received ", ws)
	passOrder(ws)
	return c.Render()
}

func passOrder(deskConn revel.ServerWebSocket) {
	for {
		order := <-deskChannel
		orderResp := TohelpDesk{Command: order.Intent, Room: 5, Data: order.Data}
		revel.INFO.Println("Request to be sent for front desk: ", orderResp)

		if err := deskConn.MessageSendJSON(orderResp); err != nil {
			revel.INFO.Println("Error sending data to desk web socket connection: ", err)
			lambdaRespChan <- false

		} else {
			lambdaRespChan <- true
		}

	}
}
