package v1

import (
	"errors"

	"github.com/cylonchau/firewalldGateway/apis"
	code "github.com/cylonchau/firewalldGateway/server/apis"
	"github.com/cylonchau/firewalldGateway/utils/firewalld"

	"github.com/gin-gonic/gin"
)

type ServiceRouter struct{}

func (this *ServiceRouter) RegisterPortAPI(g *gin.RouterGroup) {
	portGroup := g.Group("/service")
	portGroup.GET("/", this.getServicesAtRuntime)
	portGroup.DELETE("/", this.deleteServicesAtRuntime)
	portGroup.POST("/", this.addServicesAtRuntime)
	portGroup.POST("/new", this.newServiceAtPermanent)
	portGroup.GET("/list", this.listServicesAtRuntime)

}

// getServicesAtRuntime ...
// @Summary getServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/get [GET]
func (this *ServiceRouter) getServicesAtRuntime(c *gin.Context) {

	var rich = &apis.Query{}
	err := c.BindQuery(rich)

	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(rich.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		apis.ConnectDbusService(c, err)
		return
	}

	services, err := dbusClient.GetServices()

	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		apis.NotFount(c, code.ErrServiceNotFount, services)
		return
	}

	apis.SuccessResponse(c, code.OK, services)
}

// addServicesAtRuntime ...
// @Summary addServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/add [GET]
func (this *ServiceRouter) addServicesAtRuntime(c *gin.Context) {

	var query = &apis.Query{}
	err := c.BindJSON(query)

	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	defer dbusClient.Destroy()
	if err != nil {
		apis.ConnectDbusService(c, err)
		return
	}
	err = dbusClient.AddService(query.Zone, query.Service, query.Timeout)
	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}
	apis.SuccessResponse(c, code.OK, nil)
}

// listServicesAtRuntime ...
// @Summary listServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [GET]
func (this *ServiceRouter) listServicesAtRuntime(c *gin.Context) {

	var query = &apis.Query{}
	err := c.Bind(query)

	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		apis.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	services, err := dbusClient.GetServices()
	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	if len(services) <= 0 {
		apis.NotFount(c, errors.New("list of available services is not found."), nil)
		return
	}

	apis.SuccessResponse(c, code.OK, services)
}

// deleteServicesAtRuntime ...
// @Summary deleteServicesAtRuntime
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [DELETE]
func (this *ServiceRouter) deleteServicesAtRuntime(c *gin.Context) {

	var query = &apis.Query{}
	err := c.Bind(query)

	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Ip)
	if err != nil {
		apis.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	if err := dbusClient.RemoveService(query.Zone, query.Service); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}
	apis.SuccessResponse(c, code.OK, nil)
}

// newServiceAtPermanent ...
// @Summary newServiceAtPermanent
// @Produce  json
// @Success 200 {object} internal.Response
// @Router /fw/v1/service/list [POST]
func (this *ServiceRouter) newServiceAtPermanent(c *gin.Context) {

	var query = &apis.ServiceSettingQuery{}

	if err := c.BindJSON(query); err != nil {
		apis.APIResponse(c, err, nil)
		return
	}

	dbusClient, err := firewalld.NewDbusClientService(query.Host)
	if err != nil {
		apis.ConnectDbusService(c, err)
		return
	}
	defer dbusClient.Destroy()

	err = dbusClient.AddNewService(query.ServiceName, query.Setting)
	if err != nil {
		apis.APIResponse(c, err, nil)
		return
	}
	apis.SuccessResponse(c, code.OK, query.Setting)
}
