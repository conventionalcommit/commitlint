const { spawnSync } = require("child_process");

const iswin = ["win32", "cygwin"].includes(process.platform);

async function install() {
  if (process.env.CI) {
    return;
  }
  const exePath = await downloadBinary();
  if (!iswin) {
    const { chmodSync } = require("fs");
    chmodSync(exePath, "755");
  }
  // run install
  spawnSync(exePath, ["install", "-f"], {
    cwd: process.env.INIT_CWD || process.cwd(),
    stdio: "inherit",
  });
}

function getDownloadURL() {
  // Detect OS
  // https://nodejs.org/api/process.html#process_process_platform
  let goOS = process.platform;
  if (iswin) {
    goOS = "windows";
  }

  const downloadOS = `${goOS.charAt(0).toUpperCase()}${goOS.slice(1)}`;

  // Detect architecture
  // https://nodejs.org/api/process.html#process_process_arch
  let arch = process.arch;
  switch (process.arch) {
    case "x64": {
      arch = "x86_64";
      break;
    }
  }
  const version = require("./package.json").version;

  return `https://github.com/conventionalcommit/commitlint/releases/download/v${version}/commitlint_${version}_${downloadOS}_${arch}.tar.gz`;
}

const path = require("path");
const tar = require("tar-stream");
const fs = require("fs");
const zlib = require("zlib");
const fetch = (...args) => import('node-fetch').then(({default: fetch}) => fetch(...args));

const download = async (url, dest) => {
  const res = await fetch(url, { redirect: "follow" });
  await new Promise((resolve, reject) => {
    const extract = tar.extract();
    extract.on("entry", function (header, stream, next) {
      if (header.name.startsWith("commitlint")) {
        var file = fs.createWriteStream(dest);
        stream.pipe(file);
        file.on("finish", function () {
          next();
          resolve();
        });
        file.on("error", (err) => {
          file.close();

          if (err.code === "EEXIST") {
            reject("File already exists");
          } else {
            fs.unlink(dest, () => {}); // Delete temp file
            reject(err.message);
          }
        });
      } else {
        stream.on("finish", function () {
          next();
        });
        stream.resume();
      }
    });

    res.body.pipe(zlib.createGunzip()).pipe(extract);
    res.body.on("error", reject);
  });
};

async function downloadBinary() {
  const downloadURL = getDownloadURL();
  const extension = iswin ? ".exe" : "";
  const fileName = `commitlint${extension}`;
  const binDir = path.join(__dirname, "bin");

  await download(downloadURL, path.join(binDir, fileName));

  return path.join(binDir, fileName);
}

// start:
install().catch((e) => {
  console.error(e);
});
