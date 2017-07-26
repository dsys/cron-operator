const HTMLWebpackPlugin = require("html-webpack-plugin");
const path = require("path");
const webpack = require("webpack");

const ENTRY_PATH = path.resolve(__dirname, "ui/index.js");
const INDEX_PATH = path.resolve(__dirname, "ui/index.html");

module.exports = {
  target: "web",
  devtool:
    process.env.NODE_ENV === "development" ? "eval-source-map" : "source-map",
  entry: ["babel-polyfill", ENTRY_PATH],
  output: {
    filename:
      process.env.NODE_ENV === "development"
        ? "assets/bundle.js"
        : "assets/bundle.[hash].js",
    publicPath: "/"
  },
  resolve: { extensions: [".js", ".jsx", ".json"] },
  module: {
    rules: [
      { test: /\.jsx?$/, loaders: ["babel-loader"], exclude: /node_modules/ },
      { test: /\.json$/, loader: "json-loader" },
      {
        test: /\.(png|jpe?g|gif|ttf|woff|eot)$/,
        loader: "file-loader?name=assets/[hash].[ext]"
      }
    ]
  },
  plugins: [
    new webpack.DefinePlugin({
      "process.env.NODE_ENV": JSON.stringify(process.env.NODE_ENV)
    }),
    new webpack.NoEmitOnErrorsPlugin(),
    new HTMLWebpackPlugin({ template: INDEX_PATH })
  ]
};
