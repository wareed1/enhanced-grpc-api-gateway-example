window.onload = function() {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    urls: [{"url":"smpl/api/orders/v1/api_orders_service.swagger.json","name":"smpl/api/orders/v1/api_orders_service.swagger.json"},{"url":"smpl/api/users/v1/api_users_service.swagger.json","name":"smpl/api/users/v1/api_users_service.swagger.json"},{"url":"smpl/orders/v1/orders_service.swagger.json","name":"smpl/orders/v1/orders_service.swagger.json"},{"url":"smpl/users/v1/users_service.swagger.json","name":"smpl/users/v1/users_service.swagger.json"}],
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
