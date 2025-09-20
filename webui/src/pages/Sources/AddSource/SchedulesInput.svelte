<script lang="ts">
  import { getSourceServiceClient } from "#/connectrpc";
  import { toDate } from "#/connectrpc/js_date";
  import { Combobox } from "bits-ui";
  import { resource, watch } from "runed";
  import { theme } from "#/store/theme";
  import cronstrue from "cronstrue";
  import PopoverInfo from "#/components/PopoverInfo.svelte";
  import IconTrash from "@lucide/svelte/icons/trash-2";
  import { Button } from "bits-ui";

  interface Props {
    value: string[];
    inputValue: string;
  }

  let { value = $bindable([]), inputValue = $bindable("") }: Props = $props();

  let cronExpressionError = $derived.by(() => {
    if (!inputValue.trim()) {
      return "";
    }
    try {
      cronstrue.toString(inputValue, {
        throwExceptionOnParseError: true,
      });
      return "";
    } catch (e) {
      return e instanceof Error ? e.message : `${e}`;
    }
  });

  const scheduleInputNextRunResource = resource(
    () => inputValue,
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

  watch(
    () => scheduleInputNextRunResource.error,
    (err) => {
      if (err) {
        cronExpressionError = err instanceof Error ? err.message : `${err}`;
      }
    },
  );

  async function getNextRun(pattern: string, signal?: AbortSignal) {
    cronstrue.toString(pattern);
    const client = await getSourceServiceClient({ signal });
    return client.getCronNextTime({
      cronExpression: pattern,
    });
  }

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
    value: string;
    label: string;
  };

  function toPatternOption(value: string): PatternOption {
    let label;
    try {
      label = cronstrue.toString(value, {
        throwExceptionOnParseError: true,
        verbose: true,
        use24HourTimeFormat: true,
      });
    } catch (e) {
      label = `Invalid pattern: ${e instanceof Error ? e.message : e}`;
    }
    return { value, label };
  }

  function genDropdownDays(): PatternOption[] {
    const hours = Array.from({ length: 24 }).keys();
    const dayOfWeeks = ["SUN", "MON", "TUE", "WED", "THU", "FRI", "SAT"];
    const patterns: PatternOption[] = [];

    for (const hour of hours) {
      for (const day of dayOfWeeks) {
        patterns.push(toPatternOption(`0 ${hour} * * ${day}`));
      }
      patterns.push(toPatternOption(`0 ${hour} * * *`));
    }
    return patterns;
  }

  const patternDropdownList: PatternOption[] = [
    toPatternOption("* * * * *"),
    toPatternOption("*/5 * * * *"),
    toPatternOption("*/15 * * * *"),
    toPatternOption("*/30 * * * *"),
    toPatternOption("0 * * * *"),
    toPatternOption("0 0 * * MON,TUE,WED,THU,FRI"),
  ].concat(genDropdownDays());

  const filteredPatternDropdownList = $derived.by(() => {
    const text = inputValue.trim().toLowerCase();
    if (!text) {
      return patternDropdownList;
    }
    const filtered = patternDropdownList.filter(
      (opt) =>
        opt.value.toLowerCase().includes(text) ||
        opt.label.toLowerCase().includes(text),
    );
    return filtered;
  });

  let showDropdown = $state(false);

  let inputField: HTMLInputElement | null = $state(null);

  let isMobile = $state(window.innerWidth < 640);

  function handleAddSchedule() {
    if (value.includes(inputValue.trim())) {
      cronExpressionError = "This pattern is already added.";
      return;
    }
    value.push(inputValue.trim());
    inputValue = "";
  }
  function handleOnRemoveSchedule(_: ScheduleEntry, index: number) {
    value.splice(index, 1);
  }
</script>

<svelte:window onresize={() => (isMobile = window.innerWidth < 640)} />

{#snippet helpText()}
  <PopoverInfo title="Schedule Help">
    <div class="prose text-wrap">
      <p>
        Schedule uses standard 5-element Cron expression format: <span
          class="font-bold"
        >
          * * * * *
        </span>
        . Claw does not support non-standard extensions like "year" or "seconds"
        field.
      </p>
      <p>
        You can use the dropdown when typing to select common cron expression
        patterns or type your own pattern. The dropdown will filter the options
        as you type. <span class="font-bold italic">
          However, this is not an exhaustive list of cron expression patterns.
          Cron expressions can be very complex and flexible, so the dropdown
          only provides some common examples.
        </span>
        .
      </p>
      <p>
        For more advanced scheduling, you can use Google or other search engine
        using keywords like <a
          href="https://www.google.com/search?q=cron+expression+every+3+hours+on+Weekdays&sourceid=chrome&ie=UTF-8"
          class="font-bold"
        >
          cron expression every 3 hours on Weekdays
        </a>
        to find the right pattern for your needs.
      </p>
      <p>
        If you leave the field empty, the schedule will be removed and the
        source will not be scheduled to run automatically.
      </p>
      <p>
        For more information about Cron expression, please visit
        <a
          href="https://en.wikipedia.org/wiki/Cron"
          target="_blank"
          rel="noopener noreferrer"
        >
          Cron - Wikipedia
        </a>
        .
      </p>
    </div>
  </PopoverInfo>
{/snippet}

<Combobox.Root
  type="single"
  {inputValue}
  onValueChange={(value) => (inputValue = value.trim())}
  open={showDropdown}
>
  <fieldset class="fieldset w-full">
    <legend class="fieldset-legend">
      <span>Schedule</span>
      {@render helpText()}
    </legend>
    <div class="flex gap-2 items-center">
      <Combobox.Input
        placeholder="e.g. 0 0 * * FRI ({toPatternOption('0 0 * * FRI').label})"
        oninput={(e) => (inputValue = e.currentTarget.value)}
        onfocus={() => (showDropdown = true)}
        onblur={() => (showDropdown = false)}
        class="input w-full"
        bind:ref={inputField}
      />
      <button
        type="button"
        class="btn btn-primary"
        disabled={!inputValue || !!cronExpressionError}
        onclick={handleAddSchedule}
      >
        + Add
      </button>
    </div>
    {#if cronExpressionError}
      <div class="alert alert-error alert-soft">
        {cronExpressionError}
      </div>
    {:else if scheduleInputNextRunResource.loading}
      <div class="alert alert-warning alert-soft">
        <div class="loading loading-spinner"></div>
        <span>Validating cron exppression...</span>
      </div>
    {:else if scheduleInputNextRunResource.current}
      {@const data = scheduleInputNextRunResource.current}
      {@const intl = new Intl.DateTimeFormat().resolvedOptions()}
      {@const ts = toDate(data.nextTime!)}
      {@const nextRunLocal = ts.toLocaleString()}
      {@const nextRunServer = ts.toLocaleString(intl.locale, {
        timeZone: data.zone,
      })}
      <div class="alert alert-success alert-soft flex flex-col">
        <span>
          Next Run (local time): {nextRunLocal}
        </span>
        <span>Server Next Run ({data.zone}): {nextRunServer}</span>
      </div>
    {:else}
      <p class="label text-wrap">
        Schedule to extract images automatically at certain times using Cron
        expression pattern. Only standard 5 Cron expression elements are
        supported. Filling this field will enable the "Add" button (if
        expression is valid), but disables the "Save" button. To re-enable the
        "Save" button, keep this field empty.
      </p>
    {/if}
    {#await schedules}
      <div class="alert alert-warning alert-soft">
        <div class="loading loading-spinner"></div>
        <span>Validating cron entries...</span>
      </div>
    {:then entries}
      {#if entries.length}
        <table class="table">
          <thead>
            <tr>
              <th></th>
              <th>Pattern</th>
              <th>Server Next Run ({entries[0].zone})</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            {#each entries as entry, i (entry.pattern)}
              <tr>
                <th>{i + 1}</th>
                <td>{entry.pattern}</td>
                {#if entry.error}
                  <td>
                    <div class="alert alert-error alert-soft">
                      {entry.error}
                    </div>
                  </td>
                {:else if entry.nextRun}
                  {@const intl = new Intl.DateTimeFormat().resolvedOptions()}
                  {@const ts = entry.nextRun!}
                  {@const nextRunServer = ts.toLocaleString(intl.locale, {
                    timeZone: entry.zone,
                  })}
                  <td>
                    {nextRunServer}
                  </td>
                {/if}
                <th>
                  <Button.Root
                    onclick={() => handleOnRemoveSchedule(entry, i)}
                    class="btn w-full"
                  >
                    <IconTrash />
                  </Button.Root>
                </th>
              </tr>
            {/each}
          </tbody>
        </table>
      {/if}
    {/await}
  </fieldset>
  <Combobox.Portal>
    <Combobox.Content
      data-theme={$theme}
      class="z-[9999]"
      side="top"
      sideOffset={10}
    >
      <Combobox.Viewport
        class="menu bg-base-200 rounded-sm"
        style={{
          width: inputField ? `${inputField.offsetWidth}px` : "auto",
          "box-shadow": "0 -0.5rem 1.5rem rgba(0, 0, 0, 0.1)",
        }}
      >
        <div
          class="flex flex-col overflow-auto"
          style={`width: ${inputField ? `${inputField.offsetWidth}px` : "auto"}; max-height: ${isMobile ? "12.5rem" : "20rem"}`}
        >
          {#each filteredPatternDropdownList as item, i (item.value)}
            {#if i > 0}
              <div class="divider py-0 my-0"></div>
            {/if}
            <Combobox.Item
              value={item.value}
              label={item.value}
              class="btn btn-ghost flex-col items-start justify-center gap-0 py-8"
            >
              <span>{item.value}</span>
              <span class="font-normal text-left">{item.label}</span>
            </Combobox.Item>
          {:else}
            <Combobox.Item
              value={inputValue}
              label={inputValue}
              class="btn btn-ghost flex-col items-start justify-center gap-0 py-8"
            >
              <span>{inputValue}</span>
              <span class="font-normal text-left">
                {toPatternOption(inputValue).label}
              </span>
            </Combobox.Item>
          {/each}
        </div>
      </Combobox.Viewport>
    </Combobox.Content>
  </Combobox.Portal>
</Combobox.Root>
