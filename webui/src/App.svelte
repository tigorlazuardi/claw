<script lang="ts">
  import Sidebar from "./lib/Sidebar.svelte";
  import Routing from "./Routing.svelte";
  import { RouterView } from "@dvcol/svelte-simple-router/components";
</script>

<Routing>
  <main class="app">
    <Sidebar />
    <div class="content">
      <RouterView />
    </div>
  </main>
</Routing>

<style>
  .app {
    display: flex;
    height: 100vh;
    width: 100vw;
    background-color: #1a1a1a;
    color: #ffffff;
    margin: 0;
    padding: 0;
  }

  .content {
    flex: 1;
    overflow: auto;
  }
</style>
