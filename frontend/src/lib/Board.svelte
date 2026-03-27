<script>
  import { draggable, droppable } from '@thisux/sveltednd';
  import TaskCard from './TaskCard.svelte';

  let { project, onTaskClick, onTaskMove, filterLabelId = '' } = $props();

  function tasksForColumn(columnId) {
    return (project.tasks || [])
      .filter(t => t.columnId === columnId)
      .filter(t => !filterLabelId || t.labelId === filterLabelId)
      .sort((a, b) => a.position - b.position);
  }

  function handleDrop(state) {
    const { sourceContainer, targetContainer, draggedItem } = state;

    if (!draggedItem || !targetContainer) return;

    // Find the target column's tasks to determine position
    const targetTasks = tasksForColumn(targetContainer);
    const position = targetTasks.length; // append to end

    // Only call move if something actually changed
    if (draggedItem.columnId !== targetContainer || true) {
      onTaskMove?.(draggedItem.id, targetContainer, position);
    }
  }
</script>

<div class="board">
  {#each project.columns as column (column.id)}
    {@const columnTasks = tasksForColumn(column.id)}
    <div class="column">
      <div class="column-header">
        <span class="column-name">{column.name}</span>
        <span class="column-count">{columnTasks.length}</span>
      </div>
      <div
        class="column-body"
        use:droppable={{ container: column.id, callbacks: { onDrop: handleDrop } }}
      >
        {#each columnTasks as task (task.id)}
          <div
            class="card-wrapper"
            use:draggable={{ container: column.id, dragData: task }}
          >
            <TaskCard {task} labels={project.labels} onclick={onTaskClick} />
          </div>
        {/each}
        {#if columnTasks.length === 0}
          <div class="empty-zone"></div>
        {/if}
      </div>
    </div>
  {/each}
</div>

<style>
  .board {
    display: flex;
    gap: 12px;
    padding: 16px;
    min-height: calc(100vh - 50px);
    overflow-x: auto;
  }

  .column {
    flex: 0 0 260px;
    background: #e8e8e8;
    border-radius: 6px;
    display: flex;
    flex-direction: column;
    max-height: calc(100vh - 82px);
  }

  .column-header {
    padding: 10px 12px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .column-name {
    font-weight: 600;
    font-size: 0.875rem;
    color: #333;
  }

  .column-count {
    font-size: 0.75rem;
    color: #888;
    background: #d0d0d0;
    padding: 1px 6px;
    border-radius: 10px;
  }

  .column-body {
    padding: 4px 8px 8px;
    flex: 1;
    overflow-y: auto;
    display: flex;
    flex-direction: column;
    gap: 6px;
    min-height: 60px;
  }

  .card-wrapper {
    cursor: grab;
  }

  .card-wrapper:active {
    cursor: grabbing;
  }

  :global(.card-wrapper[data-dragging="true"]) {
    opacity: 0.5;
  }

  .empty-zone {
    min-height: 40px;
  }
</style>
