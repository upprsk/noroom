const connectedSpan = document.getElementById("connected-span");
const connectedErrorDiv = document.getElementById("connected-error-div");
const refreshPodsButton = document.getElementById("refresh-pods-button");
const podListUl = document.getElementById("pod-list-ul");
const podOutputCode = document.getElementById("pod-output-code");
const createPodInput = document.getElementById("create-pod-input");
const createPodButton = document.getElementById("create-pod-button");
const createPodP = document.getElementById("create-pod-p");
const attachPodInput = document.getElementById("attach-pod-input");
const attachPodButton = document.getElementById("attach-pod-button");
const attachPodP = document.getElementById("attach-pod-p");
const sendPodForm = document.getElementById("send-pod-form");
const sendPodInput = document.getElementById("send-pod-input");
const sendPodButton = document.getElementById("send-pod-button");
const podCtrlCButton = document.getElementById("pod-ctrl+c-button");

connectedErrorDiv.innerHTML = ``;

/**
 * @param ws {WebSocket}
 */
const makeClient = (ws) => {
  return {
    nextId: 1,
    live: new Map(),
    ws,
    /**
     * @param obj {{method: string; id: number}}
     */
    send(obj, expectResponse = true) {
      if (!obj.id) {
        obj.id = this.getId();
      }

      this.ws.send(JSON.stringify(obj));
      if (expectResponse) {
        return new Promise((resolve, reject) => {
          this.live.set(obj.id, [resolve, reject]);
        });
      }
    },
    /**
     * @param data {string}
     */
    recv(data) {
      for (const line of data.split("\n")) {
        const obj = JSON.parse(line);
        if (obj.id) {
          const req = this.live.get(obj.id);
          this.live.delete(obj.id);

          if (req) {
            const [resolve, reject] = req;

            if (obj.err) reject(obj);
            else resolve(obj);
          } else {
            console.error("message not found:", obj, obj);
          }
        } else {
          this.onRecv(obj);
        }
      }
    },
    onRecv() {},
    getId() {
      return this.nextId++;
    },
  };
};

const ws = new WebSocket(`ws://localhost:8080/ws`);
window.onbeforeunload = () => {
  connectedSpan.innerText = "no";
  ws.close();
};

let c = makeClient(ws);

ws.onclose = () => {
  connectedSpan.innerText = "no";
};

ws.onerror = (e) => {
  console.error("error", e);

  connectedErrorDiv.innerHTML = `<p>Error in connection: ${e}</p>`;
};

ws.onopen = () => {
  connectedSpan.innerText = "connecting...";
  connectedErrorDiv.innerText = "";

  ws.onmessage = (e) => {
    c.recv(e.data);
  };

  c.onRecv = (obj) => {
    if (obj.name === "podOut") {
      const txt = atob(obj.body.data);
      podOutputCode.innerText = podOutputCode.innerText + txt;
    }
  };

  c.send({ method: "open" }, false);
  connectedSpan.innerText = "connected";

  (async () => {
    await updatePods();

    // await c.send({
    //   method: "uploadToPod",
    //   args: ["bob", "/home/hello.txt", btoa("hello there!")],
    // });
  })();
};

const updatePods = async () => {
  const res = await c.send({ method: "listPods" });
  console.log(res.body.pods);

  const one = ({ names, state, image, status }) =>
    `<li><b>${names.join(", ")}</b> - ${image} - ${state}: ${status}</li>`;
  podListUl.innerHTML = res.body.pods.map(one);

  return res;
};

/**
 * @param name {string}
 */
const createPod = async (name) => {
  createPodP.innerText = "";

  try {
    const res = await c.send({ method: "createPod", args: [name] });
    createPodP.innerText = `Created pod with id: ${res.body.podId}`;
  } catch (e) {
    console.error(e);
    createPodP.innerText = `Failed to create pod with name "${name}": ${e.err}`;
  }
};

/**
 * @param name {string}
 */
const attachPod = async (name) => {
  attachPodP.innerText = "";

  try {
    const res = await c.send({ method: "attachToPod", args: [name] });
    attachPodP.innerText = `attached to pod with id: ${res.body.podId}`;
  } catch (e) {
    console.error(e);
    attachPodP.innerText = `Failed to attach to pod with name "${name}": ${e.err}`;
  }
};

/**
 * @param cmd {string}
 */
const sendPod = async (cmd) => {
  try {
    await c.send({ method: "sendToPod", args: [btoa(cmd)] });
  } catch (e) {
    console.error(e);
  }
};

refreshPodsButton.onclick = updatePods;
createPodButton.onclick = async () => {
  await createPod(createPodInput.value);
  createPodInput.value = "";
};
attachPodButton.onclick = async () => {
  await attachPod(attachPodInput.value);
  attachPodInput.value = "";
};
sendPodForm.onsubmit = async (e) => {
  e.preventDefault();

  await sendPod(sendPodInput.value + "\n");
  sendPodInput.value = "";
  sendPodInput.focus();
};
podCtrlCButton.onclick = async () => {
  await sendPod("\x03");
};
