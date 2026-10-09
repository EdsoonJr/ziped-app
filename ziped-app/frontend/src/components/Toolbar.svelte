<script lang="ts">
  import { CreateNewArchive, OpenArchive, TestArchive, SelectArchiveFile, SelectSaveFile, SelectFilesToCompact, SelectExtractFolder, ExtractFiles, DeleteEntries, AddFiles, AddFolder, SelectFolderToCompact } from '../../wailsjs/go/services/ArchiveService.js';
  import { currentArchivePath, currentArchiveFiles, selectedFiles, isProcessing, progressPercent, isReadOnly } from '../stores/archiveStore';

  async function handleNew() {
    const dest = await SelectSaveFile();
    if (!dest) return;

    const files = await SelectFilesToCompact();
    if (!files || files.length === 0) return;

    $isProcessing = true;
    try {
        const result = await CreateNewArchive(dest, files);
        if (result.success) {
            await openArchiveFromPath(dest);
        } else {
            alert("Erro: " + result.error);
        }
    } catch (err) {
        alert("Erro fatal ao criar pacote: " + err);
    } finally {
        $isProcessing = false;
    }
  }

  async function handleOpen() {
    const filePath = await SelectArchiveFile();
    if (!filePath) return;
    await openArchiveFromPath(filePath);
  }

  async function openArchiveFromPath(path: string) {
    $isProcessing = true;
    $progressPercent = -1;
    try {
        const files = await OpenArchive(path);
        $currentArchivePath = path;
        $currentArchiveFiles = files;
        $selectedFiles = [];
    } catch (err) {
        alert("Erro ao abrir arquivo: " + err);
        $currentArchivePath = null;
        $currentArchiveFiles = [];
        $selectedFiles = [];
    } finally {
        $isProcessing = false;
    }
  }

  async function handleExtract() {
    if (!$currentArchivePath) return;

    const destFolder = await SelectExtractFolder();
    if (!destFolder) return;

    $isProcessing = true;
    $progressPercent = -1;
    try {
        const result = await ExtractFiles($currentArchivePath, destFolder, $selectedFiles);
        if (result.success) {
            alert(result.message);
        } else {
            alert("Erro: " + result.error);
        }
    } catch (err) {
        alert("Erro fatal ao extrair: " + err);
    } finally {
        $isProcessing = false;
    }
  }

  async function handleAddFiles() {
    if (!$currentArchivePath) return;
    const files = await SelectFilesToCompact();
    if (!files || files.length === 0) return;

    $isProcessing = true;
    $progressPercent = -1;
    try {
        const result = await AddFiles($currentArchivePath, files);
        if (result.success) {
            await openArchiveFromPath($currentArchivePath);
        } else {
            alert("Erro: " + result.error);
        }
    } catch (err) {
        alert("Erro fatal ao adicionar: " + err);
    } finally {
        $isProcessing = false;
    }
  }

  async function handleAddFolder() {
    if (!$currentArchivePath) return;
    const folder = await SelectFolderToCompact();
    if (!folder) return;

    $isProcessing = true;
    $progressPercent = -1;
    try {
        const result = await AddFolder($currentArchivePath, folder);
        if (result.success) {
            await openArchiveFromPath($currentArchivePath);
        } else {
            alert("Erro: " + result.error);
        }
    } catch (err) {
        alert("Erro fatal ao adicionar: " + err);
    } finally {
        $isProcessing = false;
    }
  }

  async function handleDelete() {
    if (!$currentArchivePath || $selectedFiles.length === 0) return;

    if (!confirm(`Deseja realmente excluir ${$selectedFiles.length} item(s) do pacote?`)) {
        return;
    }

    $isProcessing = true;
    try {
        const result = await DeleteEntries($currentArchivePath, $selectedFiles);
        if (result.success) {
            alert(result.message);
            await openArchiveFromPath($currentArchivePath);
        } else {
            alert("Erro: " + result.error);
        }
    } catch (err) {
        alert("Erro fatal ao excluir: " + err);
    } finally {
        $isProcessing = false;
    }
  }

  async function testBackend() {
    if (!$currentArchivePath) {
        alert("Abra um pacote primeiro para testar a integridade!");
        return;
    }
    $isProcessing = true;
    try {
      const result = await TestArchive($currentArchivePath);
      alert(result.message);
    } catch (err) {
      alert("Erro ao chamar Go: " + err);
    } finally {
      $isProcessing = false;
    }
  }
</script>

<div class="toolbar">
  <button class="btn-primary" on:click={handleNew} disabled={$isProcessing}>✨ Novo</button>
  <button class="btn" on:click={handleOpen} disabled={$isProcessing}>📂 Abrir</button>
  <div class="divider"></div>
  <button class="btn" on:click={handleAddFiles} disabled={$isProcessing || !$currentArchivePath || $isReadOnly}>
    ➕ Adicionar Arquivos
  </button>
  <button class="btn" on:click={handleAddFolder} disabled={$isProcessing || !$currentArchivePath || $isReadOnly}>
    📁 Adicionar Pasta
  </button>
  <div class="divider"></div>
  <button class="btn" on:click={handleExtract} disabled={$isProcessing || !$currentArchivePath}>
    📤 {$selectedFiles.length > 0 ? `Extrair Selecionados (${$selectedFiles.length})` : 'Extrair Tudo'}
  </button>
  <button class="btn-danger" on:click={handleDelete} disabled={$isProcessing || $selectedFiles.length === 0 || $isReadOnly}>🗑 Excluir</button>
  <div class="divider"></div>
  <button class="btn" on:click={testBackend} disabled={$isProcessing || !$currentArchivePath}>🛡 Testar</button>
  
  <div style="flex-grow: 1;"></div>
  {#if $isReadOnly}
    <div class="readonly-badge">🔒 Somente Leitura</div>
  {/if}
</div>

<style>
  .toolbar {
    display: flex;
    gap: 8px;
    padding: 12px 16px;
    background-color: var(--bg-panel);
    border-bottom: 1px solid var(--border-color);
    align-items: center;
  }
  .btn, .btn-primary, .btn-danger {
    padding: 8px 12px;
    border-radius: 6px;
    font-size: 14px;
    font-weight: 500;
    transition: all 0.2s ease;
    display: flex;
    align-items: center;
    gap: 6px;
  }
  .btn {
    background: transparent;
  }
  .btn:hover {
    background-color: rgba(255, 255, 255, 0.1);
  }
  .btn-primary {
    background-color: var(--accent);
    color: white;
  }
  .btn-primary:hover {
    background-color: var(--accent-hover);
  }
  .btn-danger {
    color: #ef4444;
  }
  .btn-danger:hover {
    background-color: rgba(239, 68, 68, 0.1);
  }
  .divider {
    width: 1px;
    height: 24px;
    background-color: var(--border-color);
    margin: 0 4px;
  }
  .btn:disabled, .btn-primary:disabled, .btn-danger:disabled {
    opacity: 0.4;
    cursor: not-allowed;
    pointer-events: none;
  }
  .readonly-badge {
    background-color: rgba(239, 68, 68, 0.15);
    color: #ef4444;
    padding: 4px 10px;
    border-radius: 4px;
    font-size: 12px;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 4px;
  }
</style>
