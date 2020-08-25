package constants

import "time"

const DEFAULT_PORT = "51820"
const DEFAULT_PORT_NUMBER = 51820
const KEEY_ALIVE_SECONDS = 90 * time.Second

type LinkType uint8

const (
	LINK_TYPE_DIRECT LinkType = 0
	LINK_TYPE_OBFS            = 1
)
