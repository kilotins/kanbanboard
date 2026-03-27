<script>
  import { getSetupStatus, getAppTitle, getMe, logout as apiLogout, listProjects, getProject } from './lib/api.js';
  import Onboarding from './lib/Onboarding.svelte';
  import Login from './lib/Login.svelte';
  import ProjectDropdown from './lib/ProjectDropdown.svelte';
  import CreateProjectModal from './lib/CreateProjectModal.svelte';

  let loading = $state(true);
  let setupRequired = $state(false);
  let appTitle = $state('Kanban Board');
  let currentUser = $state(null);
  let projects = $state([]);
  let currentProject = $state(null);
  let showCreateProject = $state(false);

  async function checkStatus() {
    loading = true;
    try {
      const [status, titleData] = await Promise.all([
        getSetupStatus(),
        getAppTitle().catch(() => ({ title: 'Kanban Board' })),
      ]);
      setupRequired = status.setupRequired;
      appTitle = titleData.title;

      if (!setupRequired) {
        try {
          currentUser = await getMe();
          await loadProjects();
        } catch {
          currentUser = null;
        }
      }
    } catch {
      // API unreachable
    } finally {
      loading = false;
    }
  }

  async function loadProjects() {
    projects = await listProjects();
    if (projects.length > 0 && !currentProject) {
      await selectProject(projects[0]);
    }
  }

  async function selectProject(project) {
    currentProject = await getProject(project.id);
  }

  function handleSetupComplete() {
    checkStatus();
  }

  async function handleLogin(user) {
    currentUser = user;
    await loadProjects();
  }

  async function handleLogout() {
    await apiLogout();
    currentUser = null;
    projects = [];
    currentProject = null;
  }

  async function handleProjectCreated(project) {
    showCreateProject = false;
    projects = [...projects, project];
    currentProject = project;
  }

  $effect(() => {
    checkStatus();
  });
</script>

{#if loading}
  <div class="center">
    <p>Loading...</p>
  </div>
{:else if setupRequired}
  <Onboarding onComplete={handleSetupComplete} />
{:else if !currentUser}
  <Login {appTitle} onLogin={handleLogin} />
{:else}
  <div class="app">
    <header>
      <div class="header-left">
        <ProjectDropdown
          {projects}
          {currentProject}
          onSelect={selectProject}
          onCreateNew={() => showCreateProject = true}
        />
      </div>
      <div class="header-right">
        <span class="user-name">{currentUser.name}</span>
        <button class="sign-out" onclick={handleLogout}>Sign Out</button>
      </div>
    </header>

    <main>
      {#if projects.length === 0}
        <div class="center empty-state">
          <h2>Welcome to {appTitle}</h2>
          <p>Create your first project to get started.</p>
          <button class="create-btn" onclick={() => showCreateProject = true}>
            Create Project
          </button>
        </div>
      {:else if currentProject}
        <div class="board">
          {#each currentProject.columns as column}
            <div class="column">
              <div class="column-header">{column.name}</div>
              <div class="column-body">
                <p class="placeholder">Tasks coming in Phase 2.2</p>
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </main>
  </div>

  {#if showCreateProject}
    <CreateProjectModal
      onCreated={handleProjectCreated}
      onCancel={() => showCreateProject = false}
    />
  {/if}
{/if}

<style>
  .center {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 100vh;
  }

  .app {
    min-height: 100vh;
    display: flex;
    flex-direction: column;
    background: #f5f5f5;
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 16px;
    border-bottom: 1px solid #e0e0e0;
    background: #fff;
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-name {
    color: #555;
    font-size: 0.875rem;
  }

  .sign-out {
    padding: 6px 12px;
    background: none;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
    color: #555;
  }

  .sign-out:hover {
    background: #f5f5f5;
  }

  main {
    flex: 1;
    overflow-x: auto;
  }

  .empty-state {
    min-height: calc(100vh - 50px);
  }

  .empty-state h2 {
    color: #333;
    margin: 0 0 8px;
  }

  .empty-state p {
    color: #666;
    margin: 0 0 24px;
  }

  .create-btn {
    padding: 10px 24px;
    background: #4a90d9;
    color: white;
    border: none;
    border-radius: 4px;
    font-size: 1rem;
    cursor: pointer;
  }

  .create-btn:hover {
    background: #357abd;
  }

  .board {
    display: flex;
    gap: 12px;
    padding: 16px;
    min-height: calc(100vh - 50px);
  }

  .column {
    flex: 0 0 260px;
    background: #e8e8e8;
    border-radius: 6px;
    display: flex;
    flex-direction: column;
  }

  .column-header {
    padding: 10px 12px;
    font-weight: 600;
    font-size: 0.875rem;
    color: #333;
  }

  .column-body {
    padding: 8px;
    flex: 1;
  }

  .placeholder {
    color: #999;
    font-size: 0.8rem;
    text-align: center;
    padding: 16px 0;
  }

  p {
    color: #666;
  }
</style>
