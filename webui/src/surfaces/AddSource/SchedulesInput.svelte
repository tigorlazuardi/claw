<script lang="ts">
  import { getSourceServiceClient } from "../../connectrpc";
  import { toDate } from "../../connectrpc/js_date";
  import { parseCronExpression } from "cron-schedule";
  import { Popover, Combobox } from "bits-ui";
  import IconInfo from "@lucide/svelte/icons/info";
  import { resource } from "runed";

  interface Props {
    value: string[];
  }

  let scheduleInputValue = $state("");

  const scheduleInputNextRunResource = resource(
    () => scheduleInputValue,
    async (pattern, _, { signal }) => {
      if (!pattern.trim()) {
        return;
      }
      return getNextRun(pattern, signal);
    },
    {
      debounce: 300,
    },
  );

  async function getNextRun(pattern: string, signal?: AbortSignal) {
    parseCronExpression(pattern);
    const client = await getSourceServiceClient({ signal });
    return client.getCronNextTime({
      cronExpression: pattern,
    });
  }

  const { value = $bindable([]) }: Props = $props();

  type ScheduleEntry = {
    pattern: string;
    nextRun?: Date;
    error?: string;
    zone?: string;
  };

  const schedules: Promise<ScheduleEntry[]> = $derived(
    Promise.allSettled(value.map((pattern) => getNextRun(pattern))).then(
      (results) =>
        results.map((res, idx) =>
          res.status === "fulfilled"
            ? {
                pattern: value[idx],
                nextRun: toDate(res.value.nextTime!),
                zone: res.value.zone || "UTC",
              }
            : {
                pattern: value[idx],
                error:
                  res.reason instanceof Error
                    ? res.reason.message
                    : `${res.reason}`,
              },
        ),
    ),
  );

  type PatternOption = {
    pattern: string;
    description: string;
  };

  function genDropdownDays(): PatternOption[] {
    const hours = Array.from({ length: 24 }).keys();
    const dayOfWeeks = [
      "Sunday",
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
    ];
    const patterns: PatternOption[] = [];
    for (const hour of hours) {
      for (const day of dayOfWeeks) {
        const d = day.slice(0, 3).toUpperCase();
        patterns.push({
          pattern: `0 ${hour} * * ${d}`,
          description: `At ${hour}:00 on ${day}`,
        });
      }
      patterns.push({
        pattern: `0 ${hour} * * *`,
        description: `At ${hour}:00 every day`,
      });
    }
    return patterns;
  }

  const patternDropdownList: PatternOption[] = [
    {
      pattern: "* * * * *",
      description: "Runs every minute",
    },
    {
      pattern: "@minutely",
      description: "Runs every minute",
    },
    {
      pattern: "*/5 * * * *",
      description: "Runs every 5 minutes",
    },
    {
      pattern: "*/15 * * * *",
      description: "Runs every 15 minutes",
    },
    {
      pattern: "*/30 * * * *",
      description: "Runs every 30 minutes",
    },
    {
      pattern: "0 * * * *",
      description: "Runs at the start of every hour",
    },
  ].concat(genDropdownDays());

  const filteredPatternDropdownList = $derived.by(() => {
    const text = scheduleInputValue.trim().toLowerCase();
    if (!text) {
      return patternDropdownList;
    }
    return patternDropdownList.filter(
      (opt) =>
        opt.pattern.toLowerCase().includes(text) ||
        opt.description.toLowerCase().includes(text),
    );
  });

  let showDropdown = $state(false);
</script>

<fieldset class="fieldset">
  <legend class="fieldset-legend">
    <span>Schedule</span>
  </legend>
  <div class="dropdown dropdown-top">
    <input
      name="schedule"
      type="text"
      class="input w-full"
      bind:value={scheduleInputValue}
      placeholder="e.g. 0 0 * * FRI (every midnight at Friday)"
      onfocus={() => (showDropdown = true)}
      onblur={() => {
        showDropdown = false;
      }}
    />
  </div>
  {#if scheduleInputNextRunResource.error}
    <div class="alert alert-error alert-soft">
      <span>{scheduleInputNextRunResource.error.message}</span>
    </div>
  {:else if scheduleInputNextRunResource.loading}
    <div class="alert alert-warning alert-soft">
      <span class="loading loading-spinner"></span>
      <span>Validating cron expression...</span>
    </div>
  {:else if scheduleInputNextRunResource.current}
    {@const data = scheduleInputNextRunResource.current}
    {@const locale = Intl.NumberFormat().resolvedOptions().locale}
    {@const date = toDate(data.nextTime!)}
    <div class="alert alert-success alert-soft mt-2">
      <span>
        Next run time: {date.toLocaleString(locale)} (adjusted to current timezone
        {locale}). Server time: {date.toLocaleString(locale, {
          timeZone: data.zone,
        })} ({data.zone}).
      </span>
    </div>
  {:else}
    <p class="label">Schedule using Cron Expression</p>
  {/if}
</fieldset>

{#snippet dropdownList(list: PatternOption[])}
  <ul class="dropdown-content menu w-full bg-base-100 shadow-sm z-[999999]">
    {#each list as { pattern, description } (pattern)}
      <li
        class="btn flex-col items-start justify-center p-0 px-4 m-0 h-[4rem] gap-0"
      >
        <span>{pattern}</span>
        <span class="font-normal">{description}</span>
      </li>
    {/each}
  </ul>
{/snippet}

{#snippet comboInput()}
  <Combobox.Root
    type="single"
    bind:value={scheduleInputValue}
    open={showDropdown}
  >
    <fieldset class="fieldset">
      <Combobox.Input></Combobox.Input>
    </fieldset>
  </Combobox.Root>
{/snippet}

{#snippet helpText()}
  <Popover.Root>
    <Popover.Trigger class="btn btn-square btn-ghost btn-xs">
      <IconInfo />
    </Popover.Trigger>
    <Popover.Portal>
      <Popover.Content class="z-[9999] bg-transparent" data-theme="dracula"
      ></Popover.Content>
    </Popover.Portal>
  </Popover.Root>
{/snippet}
