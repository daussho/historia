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

  const token = getToken();
  if (!token) {
    console.log("no token");
    return;
  }

  const { arch, browserName, browserVersion, os } = info.platform;
  const device_name = GM_getValue(
    "device_name",
    `${os}-${arch}-${browserName}-${browserVersion}`
  );

  const host = GM_getValue("host", "");
  let titleChange = true;
  let id = "";
  while (host) {
    const title = document.title;
    if (titleChange) {
      id = await saveVisit(host, {
        title: title,
        url: document.URL,
        device_name: device_name,
      });
    }

    await new Promise((r) => setTimeout(r, 3_000));

    if (id) {
      await updateVisit(host, id);
    }

    titleChange = title !== document.title;
    console.log({ titleChange, title, "document.title": document.title });
  }

  console.log("no host");
}

async function saveVisit(host, req) {
  const res = await fetch(`${host}/api/history`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${getToken()}`,
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
      Authorization: `Bearer ${getToken()}`,
    },
  });

  const data = await res.json();
  console.log({ data });
}

function getToken() {
  return GM_getValue("token", "");
}
