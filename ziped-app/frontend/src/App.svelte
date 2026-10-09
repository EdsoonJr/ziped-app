<script lang="ts">
  import { onMount } from 'svelte';
  import Toolbar from './components/Toolbar.svelte';
  import FileList from './components/FileList.svelte';
  import { currentArchivePath, currentArchiveFiles, selectedFiles, isProcessing, progressPercent, currentProgressFile, progressStatus, isReadOnly } from './stores/archiveStore';
  import { AddFiles, SelectSaveFile, CreateNewArchive, OpenArchive } from '../wailsjs/go/services/ArchiveService.js';
  import { EventsOn } from '../wailsjs/runtime/runtime.js';

  EventsOn("progress", (data: any) => {
    if (typeof data === "number") {
      $progressPercent = data;
    } else if (data) {
      $progressPercent = data.percent || 0;
      $currentProgressFile = data.currentFile || "";
      $progressStatus = data.status || "";
    }
  });

  async function openArchiveFromPath(path: string) {
    $isProcessing = true;
    try {
      const files = await OpenArchive(path);
      $currentArchivePath = path;
      $currentArchiveFiles = files;
      $selectedFiles = [];
    } catch (err) {
      alert("Erro ao abrir arquivo: " + err);
    } finally {
      $isProcessing = false;
      $progressPercent = -1;
      $currentProgressFile = "";
      $progressStatus = "";
    }
  }

  let dragHover = false;

  async function handlePaths(paths: string[]) {
    dragHover = false;
    if ($isProcessing || paths.length === 0) return;

    if ($isReadOnly) {
      alert("O pacote atual é apenas leitura. Operação cancelada.");
      return;
    }

    setTimeout(async () => {
      if ($currentArchivePath) {
        $isProcessing = true;
        try {
          const result = await AddFiles($currentArchivePath, paths);
          if (result.success) {
            await openArchiveFromPath($currentArchivePath);
          } else {
            alert("Erro: " + result.error);
          }
        } finally {
          $isProcessing = false;
        }
      } else {
        const destPath = await SelectSaveFile();
        if (destPath) {
          $isProcessing = true;
          try {
            const result = await CreateNewArchive(destPath, paths);
            if (result.success) {
              await openArchiveFromPath(destPath);
            } else {
              alert("Erro: " + result.error);
            }
          } finally {
            $isProcessing = false;
          }
        }
      }
    }, 10);
  }

  onMount(() => {
    window.addEventListener('dragover', (e) => {
      e.preventDefault();
      dragHover = true;
    });

    window.addEventListener('dragleave', (e) => {
      e.preventDefault();
      if (e.clientX === 0 && e.clientY === 0) {
        dragHover = false;
      }
    });


    EventsOn("wails:file-drop", (x: number, y: number, paths: string[]) => {
      handlePaths(paths);
    });

    window.addEventListener('drop', (e) => {
      e.preventDefault();
      dragHover = false;
    });
  });
</script>

<main class="app-container">
  {#if dragHover}
    <div class="drag-overlay">
      <div class="drag-message">
        {$currentArchivePath ? 'Solte para adicionar ao pacote' : 'Solte para criar um novo pacote'}
      </div>
    </div>
  {/if}
  <Toolbar />
  <FileList />
  
  <div class="status-bar">
    <span style="display: flex; align-items: center; gap: 10px; flex: 1; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
      {#if $isProcessing}
        {#if $progressStatus}
          {$progressStatus}
        {:else}
          Aguarde, processando...
        {/if}

        {#if $progressPercent >= 0}
          <progress value={$progressPercent} max="100"></progress> {$progressPercent}%
        {:else}
          <progress></progress>
        {/if}

        {#if $currentProgressFile}
          <span style="margin-left: 10px; opacity: 0.8; overflow: hidden; text-overflow: ellipsis;">
            {$currentProgressFile}
          </span>
        {/if}
      {:else if $currentArchivePath}
        {$currentArchivePath}
      {:else}
        Pronto
      {/if}
    </span>
    <span style="white-space: nowrap; margin-left: 20px;">{$currentArchiveFiles.length} itens listados</span>
  </div>
</main>

<style>
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
  }
  .status-bar {
    background-color: var(--bg-panel);
    border-top: 1px solid var(--border-color);
    padding: 6px 16px;
    font-size: 12px;
    color: var(--text-muted);
    display: flex;
    justify-content: space-between;
  }
  .drag-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(27, 38, 54, 0.85);
    z-index: 1000;
    display: flex;
    justify-content: center;
    align-items: center;
    pointer-events: none;
    backdrop-filter: blur(4px);
  }
  .drag-message {
    background-color: var(--accent);
    color: white;
    padding: 20px 40px;
    border-radius: 12px;
    font-size: 24px;
    font-weight: bold;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.5);
  }
</style>
