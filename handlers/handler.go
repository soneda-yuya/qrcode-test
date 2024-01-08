package handlers

type Handler struct {
	Ticket TicketHandler
}

func NewHandler() *Handler {
	return &Handler{
		Ticket: NewTicketHandler(),
	}
}
