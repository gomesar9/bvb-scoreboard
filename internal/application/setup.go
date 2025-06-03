package application

import (
	"fmt"

	"github.com/gomesar9/bvb-scoreboard/internal/domain/service"
)

func Setup(c *service.MediaCore) error {
	lm := make(map[string]service.Listener)
	mm := make(map[string]service.MediaMaker)
	pm := make(map[string]service.Publisher)

	// Instancia Publishers
	for _, cp := range c.CoreConfig.Publishers {
		switch cp.Kind {
		case service.PublisherKindLocal:
			// TODO: Instanciar novo tipo de publisher
			pm[cp.Name] = nil
		default:
			return fmt.Errorf("Listener kind \"%s\" inválido", cp.Kind)
		}
	}

	// Instancia MediaMakers
	for _, cm := range c.CoreConfig.Makers {
		switch cm.Kind {
		case service.MediaMakerKindHTML:
			// TODO: Instanciar novo tipo de media maker
			mm[cm.Name] = nil

			// Adiciona os publishers mapeados
			for _, p := range cm.Publishers {
				mm[cm.Name].AddPublisher(pm[p])
			}
		default:
			return fmt.Errorf("MediaMaker kind \"%s\" inválido", cm.Kind)
		}
	}

	// Instancia Listeners
	for _, cl := range c.CoreConfig.Listeners {
		switch cl.Kind {
		case service.ListenerKindAPI:
			// TODO: Instanciar novo tipo de listener
			lm[cl.Name] = nil

			// Adiciona os media makers mapeados
			for _, m := range cl.Makers {
				lm[cl.Name].AddMaker(mm[m])
			}

		default:
			return fmt.Errorf("Listener kind \"%s\" inválido", cl.Kind)
		}

	}

	return nil
}
