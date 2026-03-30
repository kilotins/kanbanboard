<script>
  import { searchTasks } from './api.js';

  let { onSelect, onClose } = $props();

  let query = $state('');
  let results = $state([]);
  let searching = $state(false);
  let debounceTimer = $state(null);

  function handleInput() {
    if (debounceTimer) clearTimeout(debounceTimer);

    if (!query.trim()) {
      results = [];
      return;
    }

    debounceTimer = setTimeout(async () => {
      searching = true;
      try {
        results = await searchTasks(query.trim());
      } catch {
        results = [];
      } finally {
        searching = false;
      }
    }, 300);
  }

  function handleKeydown(e) {
    if (e.key === 'Escape') {
      onClose?.();
    }
  }

  function handleResultClick(result) {
    onSelect?.(result, result.projectId);
  }

  function autofocus(node) {
    node.focus();
  }
</script>

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div class="panel" onkeydown={handleKeydown}>
  <div class="panel-header">
    <span class="panel-title">Search Tasks</span>
    <button class="close-btn" onclick={onClose}>✕</button>
  </div>

  <div class="search-input-wrap">
    <input
      type="text"
      placeholder="Search by title or task number (e.g. KB-7)..."
      bind:value={query}
      oninput={handleInput}
      use:autofocus
    />
  </div>

  <div class="results">
    {#if searching}
      <p class="status">Searching...</p>
    {:else if query.trim() && results.length === 0}
      <p class="status">No results found</p>
    {:else}
      {#each results as result (result.id)}
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <div class="result-card" onclick={() => handleResultClick(result)}>
          <div class="result-top">
            <span class="result-ref">{result.projectTag}-{result.taskNumber}</span>
            <span class="result-project">{result.projectName}</span>
          </div>
          <span class="result-title">{result.title}</span>
        </div>
      {/each}
    {/if}
  </div>
</div>

<style>
  .panel {
    position: fixed;
    top: 0;
    right: 0;
    width: 380px;
    max-width: 90vw;
    background: white;
    box-shadow: -4px 0 16px rgba(0, 0, 0, 0.15);
    display: flex;
    flex-direction: column;
    height: 100vh;
    z-index: 200;
  }

  .panel-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid #e0e0e0;
  }

  .panel-title {
    font-weight: 600;
    font-size: 1rem;
    color: #333;
  }

  .close-btn {
    background: none;
    border: none;
    font-size: 1.2rem;
    cursor: pointer;
    color: #888;
    padding: 4px 8px;
  }

  .close-btn:hover {
    color: #333;
  }

  .search-input-wrap {
    padding: 12px 16px;
    border-bottom: 1px solid #eee;
  }

  .search-input-wrap input {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.875rem;
    box-sizing: border-box;
  }

  .search-input-wrap input:focus {
    outline: none;
    border-color: #4a90d9;
    box-shadow: 0 0 0 2px rgba(74, 144, 217, 0.2);
  }

  .results {
    flex: 1;
    overflow-y: auto;
    padding: 8px;
  }

  .status {
    text-align: center;
    color: #888;
    font-size: 0.875rem;
    padding: 24px 0;
  }

  .result-card {
    display: flex;
    flex-direction: column;
    gap: 4px;
    padding: 10px 12px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    cursor: pointer;
    margin-bottom: 6px;
  }

  .result-card:hover {
    border-color: #4a90d9;
    background: #f8fbff;
  }

  .result-top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
  }

  .result-ref {
    font-size: 0.8rem;
    font-weight: 600;
    color: #4a90d9;
  }

  .result-project {
    font-size: 0.75rem;
    color: #888;
    background: #f0f0f0;
    padding: 1px 6px;
    border-radius: 3px;
  }

  .result-title {
    font-size: 0.875rem;
    color: #333;
    line-height: 1.3;
  }
</style>
