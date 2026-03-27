<script>
  import { getSetupStatus, getAppTitle, getMe, logout as apiLogout } from './lib/api.js';
  import Onboarding from './lib/Onboarding.svelte';
  import Login from './lib/Login.svelte';

  let loading = $state(true);
  let setupRequired = $state(false);
  let appTitle = $state('Kanban Board');
  let currentUser = $state(null);

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

  function handleSetupComplete() {
    checkStatus();
  }

  function handleLogin(user) {
    currentUser = user;
  }

  async function handleLogout() {
    await apiLogout();
    currentUser = null;
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
      <h1>{appTitle}</h1>
      <div class="user-info">
        <span>{currentUser.name}</span>
        <button onclick={handleLogout}>Sign Out</button>
      </div>
    </header>
    <main class="center">
      <p>Welcome, {currentUser.name}! Board view coming in Phase 2.</p>
    </main>
  </div>
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
  }

  header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 24px;
    border-bottom: 1px solid #e0e0e0;
    background: #fff;
  }

  header h1 {
    font-size: 1.25rem;
    color: #333;
    margin: 0;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
  }

  .user-info span {
    color: #555;
    font-size: 0.875rem;
  }

  .user-info button {
    padding: 6px 12px;
    background: none;
    border: 1px solid #ccc;
    border-radius: 4px;
    font-size: 0.875rem;
    cursor: pointer;
    color: #555;
  }

  .user-info button:hover {
    background: #f5f5f5;
  }

  main {
    flex: 1;
  }

  p {
    color: #666;
  }
</style>
