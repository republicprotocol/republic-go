package hyper

import "log"

func ConsumeCertifiedBlocks(blockChIn chan Block, sb *SharedBlocks) {
	for range blockChIn {
		log.Println("Increment called in ConsumeCertifiedBlocks")
		sb.IncrementHeight()
	}
}
