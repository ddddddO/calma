const go = new Go();
let mod, instance;
WebAssembly.instantiateStreaming(fetch("main.wasm"), go.importObject).then((result) => {
    mod = result.module;
    instance = result.instance;

    console.clear();
    go.run(instance);
    instance = WebAssembly.instantiate(mod, go.importObject);
});

const callCalma = () => {
    generateCalender();
};

const clearMarkdown = () => {
  document.getElementById("in").value = "";
};

const copyToClipboard = () => {
  const calender = document.getElementById("redered_calender");
  if (calender === null) {
    return;
  }
  const clipboard = window.navigator.clipboard;
  clipboard.writeText(calender.textContent);
};

const now = () => {
  const now = new Date();
  const year = now.getFullYear();
  const month = now.getMonth() + 1;

  document.getElementById("year").value = String(year);
  document.getElementById("month").value = String(month);
};
