// ==UserScript==
// @name historia
// @namespace http://daussho.com/
// @version 0.3
// @description Historia
// @include     http://*
// @include     https://*
// @grant       GM_info
// @grant       GM_getValue
// @copyright 2024+, daussho.com
// ==/UserScript==

save();

async function save() {
  const info = GM_info;
  console.log({ info });

  await new Promise((r) => setTimeout(r, 3000));

  const token = GM_getValue("token", "");
  if (!token) {
    console.log("no token");
    return;
  }

  const { arch, browserName, browserVersion, os } = info.platform;
  const device_name = GM_getValue(
    "device_name",
    `${os}-${arch}-${browserName}-${browserVersion}`
  );
  const req = {
    title: document.title,
    url: document.URL,
    device_name: device_name,
    token: token,
  };

  console.log({ req });

  const host = GM_getValue("host", "");
  if (host) {
    const id = await saveVisit(host, req);

    while (id) {
      await new Promise((r) => setTimeout(r, 3_000));
      await updateVisit(host, id);
    }
  } else {
    console.log("no host");
  }
}

async function saveVisit(host, req) {
  const res = await fetch(`${host}/api/history`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(req),
  });

  const data = await res.json();
  console.log({ data });

  const { id } = data.data;

  return id;
}

async function updateVisit(host, id) {
  const res = await fetch(`${host}/api/history/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const data = await res.json();
  console.log({ data });
}