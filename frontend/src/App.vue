<template>
  <div class="main">
    <h1>ShortFile</h1>
    <div v-if="!file" class="content">
      <div id="drop-area">Drop file here</div>
      <div class="input-content">
        <input type="file" id="file-input" />
      </div>
    </div>
    <div v-else-if="!downloadLink && !loading">
      <div class="content">
        <h2>Hmm, Great</h2>
        <img src="./assets/happy-file.png" width="96" height="96" />
        <p>{{ file.name }}</p>
        <div>
          <button class="btn btn-upload" @click="uploadFile">Upload</button>
          <button class="btn btn-cancel" @click="reset">Cancel</button>
        </div>
      </div>
    </div>
    <div v-else-if="!downloadLink && loading">
      <div class="content">
        <h2>Uploading...</h2>
      </div>
    </div>
    <div v-else>
      <div class="content">
        <div class="download-options">
          <div>
            <span>Send Download Link via E-Mail</span>
            <div>
              <a :href="'mailto:?subject=Link to ShortFile&body=You can download the file under: http://' + downloadLink">
                <button class="btn btn-upload">Send Email</button>
              </a>
            </div>
          </div>
          <div>
            <span>Copy Download Link</span>
            <div>
              <input id="download-input" :value="'http://' + downloadLink" disabled/>
              <button id="download-copy-button" class="btn btn-upload" @click="copyUrl">Copy</button>
            </div>
          </div>
          <div>
            <span>Upload new file</span>
            <div>
              <button class="btn btn-upload" @click="reset">New Upload</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "App",
  data() {
    return {
      dropArea: null,
      fileInput: null,
      file: null,
      downloadLink: null,
      loading: false,
    };
  },
  mounted() {
    this.dropArea = document.getElementById("drop-area");
    this.fileInput = document.getElementById("file-input");

    ["dragenter", "dragover", "dragleave", "drop"].forEach((eventName) => {
      this.dropArea.addEventListener(eventName, this.preventDefaults, false);
    });

    ["dragenter", "dragover"].forEach((eventName) => {
      this.dropArea.addEventListener(eventName, this.highlight, false);
    });

    ["dragleave", "drop"].forEach((eventName) => {
      this.dropArea.addEventListener(eventName, this.unhighlight, false);
    });

    this.dropArea.addEventListener("drop", this.handleDrop, false);
    this.fileInput.addEventListener("change", this.handleFiles, false);
  },
  methods: {
    preventDefaults(e) {
      e.preventDefault();
      e.stopPropagation();
    },
    highlight(e) {
      this.dropArea.classList.add("highlight");
    },
    unhighlight(e) {
      this.dropArea.classList.remove("highlight");
    },
    handleDrop(e) {
      let dt = e.dataTransfer;
      let files = dt.files;

      this.handleFiles(files);
    },
    handleFiles(files) {
      this.file = files[0];
      //files = [...files];
      //files.forEach(this.uploadFile);
    },
    uploadFile() {
      this.loading = true;

      let url = "/u";
      let formData = new FormData();

      formData.append("file", this.file);

      this.axios
        .post(url, formData, {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        })
        .then((response) => {
          if (response.status === 200) {
            this.loading = false;
            this.downloadLink = response.data.downloadUrl;
          }
        })
        .catch((error) => {
          this.loading = false;
          console.log(error);
        });
    },
    reset() {
      this.file = null;
      this.downloadLink = null;
    },
    copyUrl() {
      let copyText = document.getElementById("download-input");
      copyText.select();
      copyText.setSelectionRange(0, 99999);

      navigator.clipboard.writeText(copyText.value);

      document.getElementById("download-copy-button").innerText = "Copied!";

      setTimeout(() => {
        document.getElementById("download-copy-button").innerText = "Copy";
      }, 5000);
    },
  },
};
</script>

<style scoped>
.main {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-start;
  padding: 50px;
  width: calc(100vw - 100px);
  height: calc(100vh - 100px);
}

.content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 50px;
}

#drop-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 400px;
  height: 200px;
  border: 2px dashed #000;
  border-radius: 10px;
  font-size: 20px;
  text-align: center;
}

#drop-area.highlight {
  border: 2px solid #000;
}

.input-content {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 20px;
  width: 400px;
  padding: 10px 0px;
  border-radius: 5px;
  background-color: white;

  box-shadow: rgba(100, 100, 111, 0.2) 0px 7px 29px 0px;
}

#file-input {
  margin-left: 5px;
}

.btn {
  width: 100px;
  height: 40px;
  border-radius: 5px;
  border: none;
  margin: 0px 10px;
  font-size: 16px;
  font-weight: bold;
  cursor: pointer;
  text-decoration: none;
}

.btn.btn-upload {
  background-color: #00b894;
  color: white;
}

.download-options {
  display: flex;
  align-items: stretch;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  width: 100%;
}

.download-options > div {
  display: flex;
  align-items: center;
  justify-content: space-around;
  flex-direction: column;
  width: 300px;
  padding: 20px 20px;
  border-radius: 10px;
  background-color: white;
  box-shadow: rgba(100, 100, 111, 0.2) 0px 7px 29px 0px;
}

.download-options > div > span {
  margin-bottom: 20px;
  font-size: 20px;
}

.download-options > div > div {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.download-options > div > div > input {
  width: calc(100% - 20px);
  height: 40px;
  border-radius: 5px;
  border: none;
  padding: 0px 10px;
  margin: 10px;
  font-size: 16px;
  font-weight: bold;
  background-color: #dfe6e9;
}
</style>
