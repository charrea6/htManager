const { createProxyMiddleware } = require('http-proxy-middleware');


module.exports = function(app) {
    app.use(
        '/api',
        createProxyMiddleware( "/api/",
            {
            target: 'http://localhost:8080',
            pathRewrite: {
                '^/api/ws': '/api/ws'
            },
            ws: true,
            changeOrigin: true,
        })
    );
};