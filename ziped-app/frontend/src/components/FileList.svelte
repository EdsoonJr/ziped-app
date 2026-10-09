<script lang="ts">
  import { currentArchiveFiles, selectedFiles } from '../stores/archiveStore';

  let currentPath = "";

  $: if ($currentArchiveFiles) {
    currentPath = "";
  }

  function formatBytes(bytes: number) {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  }

  $: visibleItems = computeVisibleItems($currentArchiveFiles, currentPath);

  function computeVisibleItems(files: any[], path: string) {
    const folders = new Map<string, boolean>();
    const result: any[] = [];
    
    if (path !== "") {
      result.push({
        isFolder: true,
        name: "..",
        fullPath: "",
        originalSize: 0,
        compressedSize: 0,
        modified: "",
        type: "Pasta acima"
      });
    }

    for (const f of files) {
      if (f.path.startsWith(path)) {
        const relative = f.path.substring(path.length);
        if (relative === "") continue;
        
        const slashIndex = relative.indexOf('/');
        if (slashIndex !== -1) {
          const folderName = relative.substring(0, slashIndex);
          if (!folders.has(folderName)) {
             folders.set(folderName, true);
             result.push({
               isFolder: true,
               name: folderName,
               fullPath: path + folderName + '/',
               originalSize: 0,
               compressedSize: 0,
               modified: "",
               type: "Pasta"
             });
          }
        } else {
          result.push({
            isFolder: false,
            name: relative,
            fullPath: f.path,
            originalSize: f.originalSize,
            compressedSize: f.compressedSize,
            modified: f.modified,
            type: f.type
          });
        }
      }
    }

    result.sort((a, b) => {
      if (a.name === "..") return -1;
      if (b.name === "..") return 1;
      if (a.isFolder && !b.isFolder) return -1;
      if (!a.isFolder && b.isFolder) return 1;
      return a.name.localeCompare(b.name);
    });

    return result;
  }

  function handleRowClick(item: any) {
    if (item.name === "..") {
      navigateUp();
    } else if (item.isFolder) {
      currentPath = item.fullPath;
    } else {
      toggleSelection(item.fullPath);
    }
  }

  function navigateUp() {
    if (!currentPath) return;
    const parts = currentPath.split('/');
    parts.pop(); 
    parts.pop(); 
    if (parts.length === 0) {
      currentPath = "";
    } else {
      currentPath = parts.join('/') + '/';
    }
  }

  function toggleSelection(filePath: string) {
    selectedFiles.update(files => {
      if (files.includes(filePath)) {
        return files.filter(f => f !== filePath);
      } else {
        return [...files, filePath];
      }
    });
  }

  function toggleAllVisible() {
    const visibleFiles = visibleItems.filter(i => !i.isFolder).map(i => i.fullPath);
    const allSelected = visibleFiles.every(f => $selectedFiles.includes(f));
    
    if (allSelected) {
      selectedFiles.update(sel => sel.filter(s => !visibleFiles.includes(s)));
    } else {
      selectedFiles.update(sel => {
        const newSel = new Set(sel);
        visibleFiles.forEach(f => newSel.add(f));
        return Array.from(newSel);
      });
    }
  }

  $: allVisibleSelected = visibleItems.filter(i => !i.isFolder).length > 0 && 
                          visibleItems.filter(i => !i.isFolder).every(i => $selectedFiles.includes(i.fullPath));

  function isItemSelected(item: any) {
    return !item.isFolder && $selectedFiles.includes(item.fullPath);
  }
</script>

<div class="file-list-container">
  {#if currentPath}
    <div class="path-bar">📂 /{currentPath}</div>
  {/if}
  <table>
    <thead>
      <tr>
        <th style="width: 40px; text-align: center;">
          <input type="checkbox" 
                 checked={allVisibleSelected}
                 on:change={toggleAllVisible} />
        </th>
        <th>Nome</th>
        <th>Tipo</th>
        <th>Tamanho Original</th>
        <th>Tamanho Compactado</th>
        <th>Modificado em</th>
      </tr>
    </thead>
    <tbody>
      {#each visibleItems as item}
        <tr class="file-row" class:selected={isItemSelected(item)} on:click={() => handleRowClick(item)}>
          <td style="text-align: center;" on:click|stopPropagation>
            {#if !item.isFolder}
              <input type="checkbox" checked={isItemSelected(item)} on:change={() => toggleSelection(item.fullPath)} />
            {/if}
          </td>
          <td>
            {#if item.isFolder}
              <span style="font-size: 1.1em; margin-right: 6px;">📁</span> {item.name}
            {:else}
              <span style="font-size: 1.1em; margin-right: 6px;">📄</span> {item.name}
            {/if}
          </td>
          <td>{item.type}</td>
          <td>{item.isFolder ? '-' : formatBytes(item.originalSize)}</td>
          <td>{item.isFolder ? '-' : formatBytes(item.compressedSize)}</td>
          <td>{item.isFolder ? '-' : item.modified}</td>
        </tr>
      {/each}
    </tbody>
  </table>
  
  {#if $currentArchiveFiles.length === 0}
    <div class="empty-state">
      <p>Nenhum pacote aberto ou pacote vazio.</p>
    </div>
  {/if}
</div>

<style>
  .file-list-container {
    flex: 1;
    overflow: auto;
    background-color: var(--bg-color);
    display: flex;
    flex-direction: column;
  }
  .path-bar {
    padding: 8px 16px;
    background-color: var(--bg-panel);
    border-bottom: 1px solid var(--border-color);
    font-family: monospace;
    font-size: 14px;
    color: var(--text-main);
  }
  table {
    width: 100%;
    border-collapse: collapse;
    text-align: left;
  }
  th {
    position: sticky;
    top: 0;
    background-color: var(--bg-panel);
    padding: 10px 16px;
    font-size: 13px;
    font-weight: 600;
    color: var(--text-muted);
    border-bottom: 1px solid var(--border-color);
    z-index: 10;
  }
  td {
    padding: 10px 16px;
    font-size: 14px;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-main);
  }
  .file-row {
    transition: background-color 0.15s ease;
    cursor: pointer;
  }
  .file-row:hover {
    background-color: rgba(255, 255, 255, 0.05);
  }
  .file-row.selected {
    background-color: rgba(59, 130, 246, 0.2);
  }
  .empty-state {
    display: flex;
    justify-content: center;
    align-items: center;
    flex: 1;
    color: var(--text-muted);
    font-style: italic;
  }
</style>
