import { fileURLToPath, URL } from "node:url";
import path from "node:path";

import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import fs from "fs";

import Inspect from "vite-plugin-inspect";

// element plus 样式自动按需导入
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";

import svgSprites from "rollup-plugin-svg-sprites";

// https://vitejs.dev/config/
export default defineConfig(({ command, mode }) => {
  console.log("vite.config defineConfig", command, mode);
  const env = loadEnv(mode, process.cwd(), "");
  console.log("vite.config env.NODE_ENV=", env.NODE_ENV);

  const optimizeDepsElementPlusIncludes = ["element-plus/es"];

  fs.readdirSync("node_modules/element-plus/es/components").forEach(
    (dirname) => {
      const cssPath = `node_modules/element-plus/es/components/${dirname}/style/css.mjs`;
      try {
        fs.accessSync(cssPath);
        optimizeDepsElementPlusIncludes.push(
          `element-plus/es/components/${dirname}/style/css`
        );
      } catch (err) {
        // 文件不存在时不做任何处理
      }
    }
  );

  return {
    base: "/", // 注意，必须以"/"结尾，BASE_URL配置
    define: {
      "process.env": env,
    },
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
      extensions: [".mjs", ".js", ".ts", ".jsx", ".tsx", ".json", ".vue"],
    },
    plugins: [
      vue(),
      Inspect(),
      AutoImport({
        resolvers: [ElementPlusResolver()],
      }),
      Components({
        resolvers: [ElementPlusResolver()],
      }),
      svgSprites({
        vueComponent: true,
        exclude: ["node_modules/**"],
        symbolId(filePath) {
          const filename = path.basename(filePath);
          return "icon-" + filename.substring(0, filename.lastIndexOf("."));
        },
      }),
    ],
    server: {
      host: "0.0.0.0",
      port: 8001,
      logLevel: "error",
      proxy: {
        "/api": {
          target: env.VITE_API_URL,
          changeOrigin: true,
          pathRewrite: {
            "^/api": "api",
          },
        },
      },
    },
  };
});
