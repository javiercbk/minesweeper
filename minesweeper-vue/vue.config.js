module.exports = {
  devServer: {
    proxy: {
      "/api": {
        target: "http://localhost:3000/",
        ws: true,
        secure: false
      }
    }
  }
};
