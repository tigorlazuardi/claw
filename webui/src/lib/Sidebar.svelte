<script lang="ts">
  import { link, active } from "@dvcol/svelte-simple-router/action";
  import { useRouter } from "@dvcol/svelte-simple-router";

  const router = useRouter();

  const menuItems = [
    {
      id: "home",
      label: "Home",
      icon: "🏠",
    },
    { id: "images", label: "Images", icon: "🖼️" },
    { id: "devices", label: "Devices", icon: "📱" },
    { id: "sources", label: "Sources", icon: "📡" },
    { id: "jobs", label: "Jobs", icon: "⚙️" },
  ];
</script>

<aside class="sidebar">
  <div class="logo">
    <h2>Claw</h2>
  </div>

  <nav class="nav">
    {#each menuItems as item}
      <a
        href={item.id === "home" ? "/" : item.id}
        class={{
          "nav-item": true,
          active: router.id === item.id,
        }}
        use:link
      >
        <span class="icon">{item.icon}</span>
        <span class="label">{item.label}</span>
      </a>
    {/each}
  </nav>
</aside>

<style>
  .sidebar {
    width: 200px;
    min-width: 200px;
    background-color: #0f0f0f;
    border-right: 1px solid #333;
    display: flex;
    flex-direction: column;
    height: 100vh;
    flex-shrink: 0;
  }

  .logo {
    padding: 1.5rem 1rem;
    border-bottom: 1px solid #333;
  }

  .logo h2 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 700;
    color: #ffffff;
  }

  .nav {
    padding: 1rem 0;
    flex: 1;
  }

  .nav-item {
    display: flex;
    align-items: center;
    width: 100%;
    padding: 0.75rem 1rem;
    color: #aaa;
    text-decoration: none;
    cursor: pointer;
    transition: all 0.2s ease;
    font-size: 0.9rem;
  }

  .nav-item:hover {
    background-color: #1f1f1f;
    color: #ffffff;
  }

  .nav-item.active {
    background-color: #2a2a2a;
    color: #ffffff;
    border-right: 3px solid #646cff;
  }

  .icon {
    margin-right: 0.75rem;
    font-size: 1.1rem;
  }

  .label {
    font-weight: 500;
  }
</style>

