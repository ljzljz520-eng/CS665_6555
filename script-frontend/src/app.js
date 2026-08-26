const state = { scripts: [], selected: null };
const views = ["list", "edit", "preview"];

function show(name) {
  for (const view of views) document.getElementById(view).hidden = view !== name;
}

async function request(path, options = {}) {
  const response = await fetch(path, { headers: { "Content-Type": "application/json" }, ...options });
  const body = await response.json();
  if (!response.ok) throw new Error(body.error || "请求失败");
  return body;
}

async function loadScripts() {
  state.scripts = await request("/api/scripts");
  const host = document.getElementById("scripts");
  host.replaceChildren(...state.scripts.map(script => {
    const row = document.createElement("button");
    row.className = "script-row";
    row.innerHTML = `<strong>${script.title}</strong><span>${script.genre || "未分类"}</span><small>${script.status}</small>`;
    row.onclick = async () => { state.selected = await request(`/api/scripts/${script.id}`); fillEditor(); show("edit"); };
    return row;
  }));
}

function fillEditor() {
  document.getElementById("title").value = state.selected?.script.title || "";
  document.getElementById("logline").value = state.selected?.script.logline || "";
  document.getElementById("genre").value = state.selected?.script.genre || "";
  document.getElementById("scenes").textContent = (state.selected?.scenes || []).map(scene => `${scene.position}. ${scene.heading}`).join("\n");
}

document.querySelectorAll("[data-view]").forEach(button => button.onclick = () => show(button.dataset.view));
document.getElementById("new-script").onclick = () => { state.selected = null; fillEditor(); show("edit"); };
document.getElementById("save-script").onclick = async () => {
  if (state.selected) return;
  const requestKey = `editor-${document.getElementById("title").value}`;
  const script = await request("/api/scripts", { method: "POST", body: JSON.stringify({ requestKey, title: document.getElementById("title").value, logline: document.getElementById("logline").value, genre: document.getElementById("genre").value }) });
  state.selected = await request(`/api/scripts/${script.id}`);
  await loadScripts();
};
document.getElementById("add-scene").onclick = async () => {
  if (!state.selected) return;
  await request(`/api/scripts/${state.selected.script.id}/scenes`, { method: "POST", body: JSON.stringify({ heading: "INT. NEW SCENE", synopsis: "待完善", location: "LOCATION", timeOfDay: "DAY" }) });
  state.selected = await request(`/api/scripts/${state.selected.script.id}`);
  fillEditor();
};
document.getElementById("refresh-preview").onclick = async () => {
  if (!state.selected) return;
  const documentView = await request(`/api/scripts/${state.selected.script.id}/preview`);
  document.getElementById("reading").textContent = documentView.text;
};

loadScripts().catch(error => { document.getElementById("scripts").textContent = error.message; });
