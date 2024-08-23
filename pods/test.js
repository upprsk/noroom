const ws = new WebSocket(`ws://localhost:8080/ws`);

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
    recv(data) {
      const obj = JSON.parse(data);
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
    },
    onRecv() {},
    getId() {
      return this.nextId++;
    },
  };
};

ws.onclose = () => {
  console.log("close");
};

ws.onerror = (e) => {
  console.error("error", e);
};

ws.onopen = () => {
  console.log("open");

  const c = makeClient(ws);
  ws.onmessage = (e) => {
    console.log(e.data);
    c.recv(e.data);
  };

  c.onRecv = (obj) => {
    console.log("got event:", obj);
  };

  c.send({ method: "open" }, false);

  (async () => {
    const res = await c.send({ method: "listPods" });
    console.log("listPods:", res);
  })();
};
