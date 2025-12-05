const form = document.getElementById("exec-form");
const languageSelect = document.getElementById("language");
const outputArea = document.getElementById("output-area");
const editorContainer = document.getElementById("editor");
const problemSelect = document.getElementById("problem-select");
const problemDescription = document.getElementById("problem-description");
const problemConstraints = document.getElementById("problem-constraints");
const submitProblemBtn = document.getElementById("submit-problem");
const testResultsList = document.getElementById("test-results");
const verdictEl = document.getElementById("verdict");
const customInputField = document.getElementById("custom-input");
const runWithInputBtn = document.getElementById("run-with-input");

const apiBase =
  window.API_BASE_URL ||
  (window.location.origin && window.location.origin.startsWith("http")
    ? window.location.origin
    : "http://localhost:3000");
const normalizedBase = apiBase.replace(/\/$/, "");
const executeEndpoint = `${normalizedBase}/api/execute`;
const problemsEndpoint = `${normalizedBase}/api/problems`;
const problemSubmitEndpoint = `${normalizedBase}/api/problems/submit`;

const applyMarkedOptions = () => {
  if (!window.marked || window.marked.__optionsApplied) {
    return;
  }
  window.marked.setOptions({
    breaks: true,
    gfm: true,
    mangle: false,
    headerIds: false,
  });
  window.marked.__optionsApplied = true;
};

applyMarkedOptions();
document.addEventListener("marked-ready", () => {
  applyMarkedOptions();
  renderProblemDescription();
});

const templates = {
  go: `package main

import "fmt"

func main() {
    fmt.Println("Hello from Go!")
}`,
  python: `print("Hello from Python!")`,
  js: `console.log("Hello from JavaScript!");`,
  cpp: `#include <iostream>
using namespace std;

int main() {
    cout << "Hello from C++!" << endl;
    return 0;
}`,
  c: `#include <stdio.h>

int main() {
    printf("Hello from C!\\n");
    return 0;
}`,
  java: `import java.util.*;

public class Main {
    public static void main(String[] args) {
        System.out.println("Hello from Java!");
    }
}`,
};

const monacoLanguageMap = {
  js: "javascript",
  python: "python",
  go: "go",
  cpp: "cpp",
  c: "c",
  java: "java",
};

let editor;
let editorReady = false;
let problemsCache = [];
let activeProblem = null;
const monacoBase =
  window.MONACO_BASE_URL ||
  "https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.45.0/min";

function initMonaco() {
  if (!window.require) {
    editorContainer.textContent =
      "Failed to load Monaco loader. Check your network connection.";
    return;
  }

  window.require.config({
    paths: {
      vs: `${monacoBase}/vs`,
    },
  });

  window.require(["vs/editor/editor.main"], () => {
    editor = monaco.editor.create(editorContainer, {
      value: templates[languageSelect.value] || "",
      language: monacoLanguageMap[languageSelect.value] || "plaintext",
      theme: "vs-dark",
      automaticLayout: true,
      minimap: { enabled: false },
      fontSize: 14,
    });
    editorReady = true;
  });
}

async function loadProblems() {
  try {
    const res = await fetch(problemsEndpoint);
    if (!res.ok) {
      throw new Error("Failed to load problems");
    }
    problemsCache = await res.json();
    renderProblemOptions();
  } catch (err) {
    problemDescription.textContent = err.message;
  }
}

function renderProblemOptions() {
  problemSelect.innerHTML = "";
  problemsCache.forEach((problem) => {
    const option = document.createElement("option");
    option.value = problem.id;
    option.textContent = problem.title;
    problemSelect.appendChild(option);
  });

  if (problemsCache.length > 0) {
    setActiveProblem(problemsCache[0].id);
  }
}

function setActiveProblem(problemId) {
  activeProblem = problemsCache.find((p) => p.id === problemId) || null;
  if (!activeProblem) {
    problemDescription.textContent = "Problem not found.";
    problemConstraints.textContent = "";
    verdictEl.textContent = "";
    testResultsList.innerHTML = "";
    return;
  }

  problemSelect.value = activeProblem.id;
  renderProblemDescription();
  problemConstraints.textContent = `Input: ${activeProblem.input_format} • Constraints: ${activeProblem.constraints}`;
  verdictEl.textContent = "";
  testResultsList.innerHTML = "";
}

function renderProblemDescription() {
  if (!activeProblem) {
    return;
  }
  if (window.marked) {
    problemDescription.innerHTML = window.marked.parse(
      activeProblem.description || ""
    );
  } else {
    problemDescription.textContent = activeProblem.description || "";
  }
}

function renderTestResults(results = []) {
  testResultsList.innerHTML = "";
  results.forEach((result) => {
    const item = document.createElement("li");
    const name = document.createElement("span");
    name.className = "test-name";
    name.textContent = result.name;

    const status = document.createElement("span");
    status.className = `status-pill status-${result.status}`;
    status.textContent = result.status;

    item.appendChild(name);
    item.appendChild(status);
    testResultsList.appendChild(item);
  });
}

initMonaco();
loadProblems();

languageSelect.addEventListener("change", (event) => {
  if (!editor || !editorReady) {
    return;
  }

  const nextLang = event.target.value;
  const model = editor.getModel();
  if (model) {
    monaco.editor.setModelLanguage(
      model,
      monacoLanguageMap[nextLang] || "plaintext"
    );
  }

  if (!editor.getValue().trim()) {
    editor.setValue(templates[nextLang] || "");
  }
});

problemSelect.addEventListener("change", (event) => {
  setActiveProblem(event.target.value);
});

form.addEventListener("submit", async (event) => {
  event.preventDefault();

  if (!editor || !editorReady) {
    outputArea.textContent = "Editor is still loading. Please wait.";
    return;
  }

  const payload = {
    code: editor.getValue(),
    language: languageSelect.value,
    mode: "standard",
    input: "",
  };

  outputArea.textContent = "Running...";

  try {
    const res = await fetch(executeEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `Request failed with ${res.status}`);
    }

    const data = await res.json();
    outputArea.textContent = data.success
      ? data.output || "(no output)"
      : data.error || data.status_message || "Execution failed";
  } catch (err) {
    outputArea.textContent = `Error: ${err.message}`;
  }
});

submitProblemBtn.addEventListener("click", async () => {
  if (!editor || !editorReady) {
    verdictEl.textContent = "Editor is still loading. Please wait.";
    return;
  }
  if (!activeProblem) {
    verdictEl.textContent = "Please select a problem first.";
    return;
  }

  verdictEl.textContent = "Judging...";
  testResultsList.innerHTML = "";

  const payload = {
    problem_id: activeProblem.id,
    code: editor.getValue(),
    language: languageSelect.value,
  };

  try {
    const res = await fetch(problemSubmitEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `Request failed with ${res.status}`);
    }

    const data = await res.json();
    verdictEl.textContent = `Verdict: ${data.verdict}`;
    renderTestResults(data.results);
  } catch (err) {
    verdictEl.textContent = `Error: ${err.message}`;
  }
});

runWithInputBtn.addEventListener("click", async () => {
  if (!editor || !editorReady) {
    outputArea.textContent = "Editor is still loading. Please wait.";
    return;
  }

  const payload = {
    code: editor.getValue(),
    language: languageSelect.value,
    mode: "standard",
    input: customInputField.value || "",
  };

  outputArea.textContent = "Running with custom input...";

  try {
    const res = await fetch(executeEndpoint, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
    });

    if (!res.ok) {
      const text = await res.text();
      throw new Error(text || `Request failed with ${res.status}`);
    }

    const data = await res.json();
    outputArea.textContent = data.success
      ? data.output || "(no output)"
      : data.error || data.status_message || "Execution failed";
  } catch (err) {
    outputArea.textContent = `Error: ${err.message}`;
  }
});

