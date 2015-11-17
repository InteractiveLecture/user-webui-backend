package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/InteractiveLecture/servicecache"
	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

func createService(id string, service string, port int, address string, tags ...string) *api.AgentService {
	return &api.AgentService{
		ID:      id,
		Service: service,
		Tags:    tags,
		Port:    port,
		Address: address,
	}

}

func prepareServer(serviceName string) *httptest.Server {
	// wir schreiben einen fake server. Dieser server stellt den eigentlichen service dar,
	// an den du deine Requests weiterleiten möchtest (beispielsweise der authentication-service).
	// Falls er angesprochen wird soll er den Statuscode 418 (im a teapot) zurückgesendet.
	serviceHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	server := httptest.NewServer(serviceHandler)

	// der servicecache wird so konfiguriert, dass er immer die url des eben erstellten server zurückliefert.
	// Wir umgehen somit die periodischen Anfragen an den discovery-service.
	fakeRetriever := func(consulAddress string) (map[string]*api.AgentService, error) {
		services := make(map[string]*api.AgentService)
		serviceUrl, _ := url.Parse(server.URL)
		parts := strings.Split(serviceUrl.Host, ":")
		host := parts[0]
		port, _ := strconv.Atoi(parts[1])
		services[serviceName] = createService(serviceName+"1", serviceName, port, host, "public", "interactive-lecture")
		return services, nil
	}
	// wir configurieren den cache und bauen anschließend unseren mock-retreiver ein.
	cache, _ := servicecache.Configure("discovery:8500", 1*time.Millisecond, serviceName)
	cache.SetServiceRetriever(fakeRetriever)
	return server
}

//Alle methoden mit TestXxxx (t *testing.T) werden von go test ausgeführt.
func TestCreateProxy(t *testing.T) {
	server := prepareServer("authentication-service")
	defer server.Close()
	// nun testen wir den eigentlichen Proxy.  createProxy erzeugt einen proxy-handler.
	proxy := createProxy("authentication-service", "/oauth/token")
	testProxy := httptest.NewServer(proxy)
	defer testProxy.Close()
	//wir senden einen normalen http-request an den proxy und erwarten, dass er ihn an den fake authentication-server weiterleitet. Es muss also der statuscode "i'm a teapot" zurückkommen.
	resp, err := http.Get(testProxy.URL)
	assert.Nil(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}
