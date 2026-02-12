package store

import (
	"fmt"
	"net/url"
)

func buildWhatsAppURL(orderID string, totalAOA float64) string {
	message := fmt.Sprintf(
		"Olá 👋\n\nQuero confirmar a compra:\n\nPedido: %s\nTotal: %.0f Kz\n\nJá fiz a transferência. Segue o comprovativo.",
		orderID,
		totalAOA,
	)

	encodeMessage := url.QueryEscape(message)

	phone := "+244938252431"

	return fmt.Sprintf(
		"https://wa.me/%s?text=%s", phone, encodeMessage,
	)
}
