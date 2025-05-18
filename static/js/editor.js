class DrawingEditor {
  constructor() {
    this.canvas = document.getElementById("drawing-canvas");
    this.ctx = this.canvas.getContext("2d");
    this.scale = { x: 1, y: 1 };
    this.isDrawing = false;
    this.lastX = 0;
    this.lastY = 0;

    this.setupCanvas();
    this.setupEventListeners();
    this.calculateScale();

    window.addEventListener("resize", () => this.calculateScale());
  }

  calculateScale() {
    const rect = this.canvas.getBoundingClientRect();
    this.scale = {
      x: this.canvas.width / rect.width,
      y: this.canvas.height / rect.height,
    };
  }

  setupCanvas() {
    this.ctx.fillStyle = "#ffffff";
    this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    this.ctx.strokeStyle = "#000000";
    this.ctx.lineWidth = 5;
    this.ctx.lineCap = "round";
    this.ctx.lineJoin = "round";
  }

  setupEventListeners() {
    // Рисование
    this.canvas.addEventListener("mousedown", (e) => this.startDrawing(e));
    this.canvas.addEventListener("mousemove", (e) => this.draw(e));
    this.canvas.addEventListener("mouseup", () => this.stopDrawing());
    this.canvas.addEventListener("mouseout", () => this.stopDrawing());

    // Инструменты
    document.getElementById("color-picker").addEventListener("input", (e) => {
      this.ctx.strokeStyle = e.target.value;
      this.ctx.fillStyle = e.target.value;
      document.getElementById("current-color-hex").textContent = e.target.value;
    });

    document.getElementById("brush-size").addEventListener("change", (e) => {
      this.ctx.lineWidth = parseInt(e.target.value);
    });

    // Цветовые пресеты
    document.querySelectorAll(".color-preset").forEach((preset) => {
      preset.addEventListener("click", () => {
        const color = preset.dataset.color;
        this.ctx.strokeStyle = color;
        this.ctx.fillStyle = color;
        document.getElementById("color-picker").value = color;
        document.getElementById("current-color-hex").textContent = color;
      });
    });

    // Кнопки
    document.getElementById("clear-btn").addEventListener("click", () => {
      if (confirm("Очистить холст?")) {
        this.setupCanvas();
      }
    });

    document
      .getElementById("save-btn")
      .addEventListener("click", () => this.saveDrawing());
    document
      .getElementById("load-btn")
      .addEventListener("click", () => this.openGallery());
  }

  getCanvasPosition(e) {
    const rect = this.canvas.getBoundingClientRect();
    return {
      x: (e.clientX - rect.left) * this.scale.x,
      y: (e.clientY - rect.top) * this.scale.y,
    };
  }

  startDrawing(e) {
    this.isDrawing = true;
    const pos = this.getCanvasPosition(e);
    this.lastX = pos.x;
    this.lastY = pos.y;

    // Рисуем начальную точку
    this.ctx.beginPath();
    this.ctx.arc(
      this.lastX,
      this.lastY,
      this.ctx.lineWidth / 2,
      0,
      Math.PI * 2
    );
    this.ctx.fill();
  }

  draw(e) {
    if (!this.isDrawing) return;

    const pos = this.getCanvasPosition(e);
    const currentX = pos.x;
    const currentY = pos.y;

    this.ctx.beginPath();
    this.ctx.moveTo(this.lastX, this.lastY);
    this.ctx.lineTo(currentX, currentY);
    this.ctx.stroke();

    this.lastX = currentX;
    this.lastY = currentY;
  }

  stopDrawing() {
    this.isDrawing = false;
  }

  async saveDrawing() {
    const title = prompt("Введите название рисунка:", "Мой рисунок");
    if (!title) return;

    try {
      const imageData = this.canvas.toDataURL("image/png");
      const response = await fetch("/api/save-drawing", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, image: imageData }),
        credentials: "include",
      });

      if (!response.ok) throw new Error("Не удалось сохранить");
      alert("Рисунок сохранён!");
    } catch (error) {
      console.error("Ошибка сохранения:", error);
      alert("Ошибка: " + error.message);
    }
  }

  openGallery() {
    const modal = document.getElementById("gallery-modal");
    modal.style.display = "flex";
    this.loadGallery();
  }

  closeGallery() {
    document.getElementById("gallery-modal").style.display = "none";
  }

  async loadGallery() {
    const galleryContent = document.getElementById("gallery-content");
    galleryContent.innerHTML = '<div class="loading-spinner"></div>';

    try {
      const response = await fetch("/api/get-drawings", {
        credentials: "include", // Важно для кук
      });

      if (!response.ok) {
        const error = await response.json();
        throw new Error(error.message || "Не удалось загрузить работы");
      }

      const drawings = await response.json();
      this.renderGallery(drawings);
    } catch (error) {
      galleryContent.innerHTML = `
            <div class="error-message">
                <i class="fas fa-exclamation-triangle"></i>
                <p>${error.message}</p>
                <button class="retry-btn">Попробовать снова</button>
            </div>
        `;
      galleryContent
        .querySelector(".retry-btn")
        .addEventListener("click", () => this.loadGallery());
    }
  }

  renderGallery(drawings) {
    const galleryContent = document.getElementById("gallery-content");

    if (drawings.length === 0) {
      galleryContent.innerHTML = `
              <div class="empty-gallery">
                  <i class="fas fa-palette"></i>
                  <p>Нет сохранённых работ</p>
              </div>
          `;
      return;
    }

    galleryContent.innerHTML = "";
    const grid = document.createElement("div");
    grid.className = "drawings-grid";

    drawings.forEach((drawing) => {
      const drawingCard = document.createElement("div");
      drawingCard.className = "drawing-card";
      drawingCard.innerHTML = `
              <div class="drawing-preview">
                  <img src="${drawing.image}" alt="${drawing.title}">
                  <div class="drawing-overlay">
                      <button class="action-btn load-btn" data-image="${
                        drawing.image
                      }">
                          <i class="fas fa-edit"></i> Загрузить
                      </button>
                      <button class="action-btn delete-btn" data-id="${
                        drawing.id
                      }">
                          <i class="fas fa-trash"></i> Удалить
                      </button>
                  </div>
              </div>
              <div class="drawing-info">
                  <h3>${drawing.title}</h3>
                  <span class="drawing-date">${new Date(
                    drawing.created_at
                  ).toLocaleDateString()}</span>
              </div>
          `;

      drawingCard.querySelector(".load-btn").addEventListener("click", () => {
        this.loadDrawing(drawing.image);
        this.closeGallery();
      });

      drawingCard
        .querySelector(".delete-btn")
        .addEventListener("click", async (e) => {
          e.stopPropagation();
          if (confirm(`Удалить "${drawing.title}"?`)) {
            const success = await this.deleteDrawing(drawing.id);
            if (success) {
              drawingCard.remove();
              // Обновляем список рисунков, если он пуст
              if (document.querySelectorAll(".drawing-card").length === 1) {
                this.loadGallery();
              }
            }
          }
        });

      grid.appendChild(drawingCard);
    });

    galleryContent.appendChild(grid);

    // Поиск
    document
      .querySelector(".search-box input")
      .addEventListener("input", (e) => {
        const searchTerm = e.target.value.toLowerCase();
        document.querySelectorAll(".drawing-card").forEach((card) => {
          const title = card.querySelector("h3").textContent.toLowerCase();
          card.style.display = title.includes(searchTerm) ? "block" : "none";
        });
      });
  }

  loadDrawing(imageData) {
    const img = new Image();
    img.onload = () => {
      const currentColor = this.ctx.strokeStyle;
      const currentSize = this.ctx.lineWidth;

      this.ctx.fillStyle = "#ffffff";
      this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
      this.ctx.drawImage(img, 0, 0);

      this.ctx.strokeStyle = currentColor;
      this.ctx.fillStyle = currentColor;
      this.ctx.lineWidth = currentSize;
    };
    img.src = imageData;
  }

  async deleteDrawing(id) {
    try {
      console.log("[DELETE] Attempting to delete drawing:", id);

      const response = await fetch(`/api/delete-drawing/${id}`, {
        method: "DELETE",
        credentials: "include",
      });

      console.log("[DELETE] Response status:", response.status);

      if (!response.ok) {
        const error = await response.text();
        console.error("[DELETE] Server error:", error);
        throw new Error(error || "Delete failed");
      }

      const result = await response.json();
      console.log("[DELETE] Success:", result);

      return true;
    } catch (error) {
      console.error("[DELETE] Full error:", error);
      alert(`Ошибка удаления: ${error.message}`);
      return false;
    }
  }
}

// Инициализация редактора
document.addEventListener("DOMContentLoaded", () => {
  const editor = new DrawingEditor();

  // Обработчики модального окна
  document
    .querySelector(".modal-close")
    .addEventListener("click", () => editor.closeGallery());
  document
    .getElementById("close-gallery-btn")
    .addEventListener("click", () => editor.closeGallery());
  document.getElementById("gallery-modal").addEventListener("click", (e) => {
    if (e.target === e.currentTarget) {
      editor.closeGallery();
    }
  });
});
