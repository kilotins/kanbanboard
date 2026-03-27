<script>
  import { getSetupStatus, getAppTitle } from './lib/api.js';
  import Onboarding from './lib/Onboarding.svelte';

  let loading = $state(true);
  let setupRequired = $state(false);
  let appTitle = $state('Kanban Board');

  async function checkStatus() {
    try {
      const [status, titleData] = await Promise.all([
        getSetupStatus(),
        getAppTitle().catch(() => ({ title: 'Kanban Board' })),
      ]);
      setupRequired = status.setupRequired;
      appTitle = titleData.title;
    } catch {
      // If API is unreachable, show loading state
    } finally {
      loading = false;
    }
  }

  function handleSetupComplete() {
    checkStatus();
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
{:else}
  <div class="center">
    <h1>{appTitle}</h1>
    <p>Login screen coming in the next phase.</p>
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

  h1 {
    font-size: 2rem;
    color: #333;
    margin: 0 0 8px;
  }

  p {
    color: #666;
  }
</style>
