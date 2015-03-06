package haproxy

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

const (
	// TEMPLATE_FILE = "../configuration/templates/haproxy_config.template"
	// CONFIG_FILE   = "/tmp/vamp_lb_test.cfg"
	// EXAMPLE       = "../test/test_config1.json"
	// JSON_FILE     = "/tmp/vamp_lb_test.json"
	// PID_FILE      = "/tmp/vamp_lb_test.pid"
	ROUTE_JSON  = "../test/test_route.json"
	GROUP_JSON  = "../test/test_service.json"
	SERVER_JSON = "../test/test_server1.json"
)

func TestConfiguration_GetRoutes(t *testing.T) {

	routes := haConfig.GetRoutes()
	if routes[0].Name != "test_route_1" {
		t.Errorf("Failed to get all frontends")
	}
}

func TestConfiguration_GetRoute(t *testing.T) {

	if route, err := haConfig.GetRoute("test_route_1"); route.Name != "test_route_1" && err == nil {
		t.Errorf("Failed to get frontend")
	}

	if _, err := haConfig.GetRoute("non_existent_route"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}

}

func TestConfiguration_AddRoute(t *testing.T) {
	j, _ := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	_ = json.Unmarshal(j, &route)

	if haConfig.AddRoute(route) != nil {
		t.Errorf("Failed to add route")
	}

	if haConfig.AddRoute(route) == nil {
		t.Errorf("Adding should fail when a route already exists")
	}
}

func TestConfiguration_UpdateRoute(t *testing.T) {
	j, _ := ioutil.ReadFile(ROUTE_JSON)
	var route *Route
	if err := json.Unmarshal(j, &route); err != nil {
		t.Errorf(err.Error())
	}
	route.Protocol = "tcp"

	if err := haConfig.UpdateRoute("test_route_2", route); err != nil {
		t.Errorf(err.Error())
	}

	if route, err := haConfig.GetRoute("test_route_2"); err != nil && route.Protocol != "tcp" {
		t.Errorf("Failed to update route")
	}
}

func TestConfiguration_GetRouteServices(t *testing.T) {

	if services, err := haConfig.GetRouteServices("test_route_1"); services[0].Name != "service_a" || err != nil {
		t.Errorf("Failed to get services")
	}

	if _, err := haConfig.GetRouteServices("non_existent_service"); err == nil {
		t.Errorf("Should return nil on non existent service")
	}
}

func TestConfiguration_GetRouteService(t *testing.T) {

	if service, err := haConfig.GetRouteService("test_route_1", "service_a"); service.Name != "service_a" || err != nil {
		t.Errorf("Failed to get service")
	}

	if _, err := haConfig.GetRouteService("non_existent_route", "service_a"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}
	if _, err := haConfig.GetRouteService("test_route_1", "non_existent_service"); err == nil {
		t.Errorf("Should return nil on non existent service")
	}

}

func TestConfiguration_AddRouteServices(t *testing.T) {
	j, _ := ioutil.ReadFile(GROUP_JSON)
	var services []*Service
	_ = json.Unmarshal(j, &services)

	route := "test_route_1"

	if haConfig.AddRouteServices(route, services) != nil {
		t.Errorf("Failed to add route")
	}

	if haConfig.AddRouteServices(route, services) == nil {
		t.Errorf("Adding should fail when a service already exists")
	}

	if haConfig.AddRouteServices("non_existent_service", services) == nil {
		t.Errorf("Should return nil on non existent route")
	}

}

func TestConfiguration_UpdateRouteService(t *testing.T) {
	j, _ := ioutil.ReadFile(GROUP_JSON)
	var services []*Service
	_ = json.Unmarshal(j, &services)

	service := services[0]
	service.Weight = 1

	if err := haConfig.UpdateRouteService("test_route_1", services[0].Name, service); err != nil {
		t.Errorf(err.Error())
	}

}

func TestConfiguration_DeleteRouteService(t *testing.T) {

	route := "test_route_1"

	if err := haConfig.DeleteRouteService(route, "service_c"); err != nil {
		t.Errorf("Failed to delete route")
	}

	if haConfig.DeleteRouteService("non_existent_route", "service_a") == nil {
		t.Errorf("Should return nil on non existent route")
	}

	if haConfig.DeleteRouteService(route, "non_existent_service") == nil {
		t.Errorf("Should return nil on non existent service")
	}
}

func TestConfiguration_GetServiceServers(t *testing.T) {

	if servers, err := haConfig.GetServiceServers("test_route_1", "service_a"); err != nil {
		t.Errorf("Failed to get servers")
	} else {
		if servers[0].Name != "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd" {
			t.Errorf("Failed to get servers")
		}
	}

	if _, err := haConfig.GetServiceServers("non_existent_route", "service_a"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}

	if _, err := haConfig.GetServiceServers("test_route_1", "non_existent_service"); err == nil {
		t.Errorf("Should return nil on non existent service")
	}

}

func TestConfiguration_GetServiceServer(t *testing.T) {

	if _, err := haConfig.GetServiceServer("test_route_1", "service_a", "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"); err != nil {
		t.Errorf("Failed to get server")
	}

	if _, err := haConfig.GetServiceServer("test_route_1", "service_a", "non_existent_server"); err == nil {
		t.Errorf("Should return nil on non existent server")
	}
}

func TestConfiguration_AddServiceServer(t *testing.T) {

	route := "test_route_1"
	service := "service_a"

	j, _ := ioutil.ReadFile(GROUP_JSON)
	var server Server

	_ = json.Unmarshal(j, &server)

	if err := haConfig.AddServiceServer(route, service, &server); err != nil {
		t.Errorf(err.Error())
	}

	if err := haConfig.AddServiceServer(route, service, &server); err == nil {
		t.Errorf("Adding should fail when a server already exists")
	}

	if err := haConfig.AddServiceServer(route, "non_existent_service", &server); err == nil {
		t.Errorf("Should return error on non existent service")
	}

	if err := haConfig.AddServiceServer("non_existent_route", service, &server); err == nil {
		t.Errorf("Should return error on non existent route")
	}

	// test should be activated when "exists" checking is in the haproxy package
	// server.Name = "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"
	// err = haConfig.AddServiceServer(route,service,&server)
	// if err == nil {
	// 	t.Errorf("Should return error on trying to create an already existing server")
	// }

}

func TestConfiguration_DeleteServiceServer(t *testing.T) {

	route := "test_route_1"
	service := "service_a"
	server := "paas.55f73f0d-6087-4964-a70e-b1ca1d5b24cd"

	if err := haConfig.DeleteServiceServer(route, service, server); err != nil {
		t.Errorf("Failed to delete server")
	}

	if err := haConfig.DeleteServiceServer(route, service, "non_existent_server"); err == nil {
		t.Errorf("Should return nil on non existent server")
	}
}

func TestConfiguration_UpdateServiceServer(t *testing.T) {

	j, _ := ioutil.ReadFile(SERVER_JSON)

	var server *Server
	_ = json.Unmarshal(j, &server)
	serverToUpdate := "server_to_be_updated"
	server.Port = 1234
	routeName := "test_route_2"
	serviceName := "service_to_be_updated"

	if err := haConfig.UpdateServiceServer(routeName, serviceName, serverToUpdate, server); err != nil {
		t.Errorf(err.Error())
	}

	if server, err := haConfig.GetServiceServer(routeName, serviceName, server.Name); err != nil && server.Port != 1234 {
		t.Errorf(err.Error())
	}
}

func TestConfiguration_DeleteRoute(t *testing.T) {

	if err := haConfig.DeleteRoute("test_route_2"); err != nil {
		t.Errorf("Failed to delete route")
	}

	if err := haConfig.DeleteRoute("non_existent_route"); err == nil {
		t.Errorf("Should return nil on non existent route")
	}
}