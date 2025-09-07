<script lang="ts">
  interface Props {
    size?: string;
    color?: string;
  }

  const { size = "2rem", color = "hsl(235, 100%, 65%)" }: Props = $props();
</script>

<svg
  width={size}
  height={size}
  viewBox="0 0 24 24"
  fill="none"
  xmlns="http://www.w3.org/2000/svg"
  class="spinning-arrow"
>
  <path
    d="M12 2L12 6M12 18L12 22M22 12L18 12M6 12L2 12M19.07 4.93L16.24 7.76M7.76 16.24L4.93 19.07M19.07 19.07L16.24 16.24M7.76 7.76L4.93 4.93"
    stroke={color}
    stroke-width="2"
    stroke-linecap="round"
    stroke-linejoin="round"
  />
  <circle cx="12" cy="12" r="3" stroke={color} stroke-width="2" fill="none" />
</svg>

<style>
  .spinning-arrow {
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    from {
      transform: rotate(0deg);
    }
    to {
      transform: rotate(360deg);
    }
  }
</style>