<script>
  let { projects = [], currentProject = null, onSelect, onCreateNew } = $props();

  let open = $state(false);

  function toggle() {
    open = !open;
  }

  function select(project) {
    onSelect(project);
    open = false;
  }

  function handleCreateNew() {
    onCreateNew();
    open = false;
  }

  function handleClickOutside(e) {
    if (!e.target.closest('.dropdown')) {
      open = false;
    }
  }
</script>

<svelte:window onclick={handleClickOutside} />

<div class="dropdown">
  <button class="trigger" onclick={toggle}>
    <span class="project-name">{currentProject?.name ?? 'Select Project'}</span>
    <span class="arrow">{open ? '▲' : '▼'}</span>
  </button>

  {#if open}
    <div class="menu">
      {#each projects as project}
        <button
          class="menu-item"
          class:active={currentProject?.id === project.id}
          onclick={() => select(project)}
        >
          {project.name}
        </button>
      {/each}
      <hr />
      <button class="menu-item new-project" onclick={handleCreateNew}>
        + New Project
      </button>
    </div>
  {/if}
</div>

<style>
  .dropdown {
    position: relative;
  }

  .trigger {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 6px 12px;
    background: none;
    border: 1px solid #ccc;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.95rem;
    color: #333;
  }

  .trigger:hover {
    background: #f5f5f5;
  }

  .project-name {
    font-weight: 500;
  }

  .arrow {
    font-size: 0.7rem;
    color: #888;
  }

  .menu {
    position: absolute;
    top: 100%;
    left: 0;
    margin-top: 4px;
    min-width: 200px;
    background: white;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    z-index: 100;
  }

  .menu-item {
    display: block;
    width: 100%;
    padding: 8px 12px;
    border: none;
    background: none;
    text-align: left;
    cursor: pointer;
    font-size: 0.875rem;
    color: #333;
  }

  .menu-item:hover {
    background: #f0f4ff;
  }

  .menu-item.active {
    background: #e8f0fe;
    font-weight: 500;
  }

  .new-project {
    color: #4a90d9;
    font-weight: 500;
  }

  hr {
    margin: 4px 0;
    border: none;
    border-top: 1px solid #eee;
  }
</style>
