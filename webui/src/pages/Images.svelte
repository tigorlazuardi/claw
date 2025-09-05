<script lang="ts">
  import type { Image } from '../gen/claw/v1/image_pb'
  import MasonryLayout from '../layouts/MasonryLayout.svelte'
  import ImageModal from '../components/ImageModal.svelte'
  
  // Dummy data for demonstration
  const dummyImages: Image[] = [
    {
      id: BigInt(1),
      sourceId: BigInt(1),
      deviceIds: [BigInt(1), BigInt(2)],
      paths: ['/path/to/image1.jpg'],
      downloadUrl: 'https://example.com/image1.jpg',
      width: 1920,
      height: 1080,
      filesize: 524288,
      thumbnailPath: '/thumbnails/image1_thumb.jpg',
      imagePath: '/images/image1.jpg',
      postAuthor: 'Artist1',
      postAuthorUrl: 'https://example.com/artist1',
      postUrl: 'https://example.com/post1',
      isFavorite: false,
      tags: ['anime', 'landscape', 'digital'],
      createdAt: undefined,
      updatedAt: undefined
    },
    {
      id: BigInt(2),
      sourceId: BigInt(1),
      deviceIds: [BigInt(1)],
      paths: ['/path/to/image2.jpg'],
      downloadUrl: 'https://example.com/image2.jpg',
      width: 1920,
      height: 1080,
      filesize: 412345,
      thumbnailPath: '/thumbnails/image2_thumb.jpg',
      imagePath: '/images/image2.jpg',
      postAuthor: 'Artist2',
      postAuthorUrl: 'https://example.com/artist2',
      postUrl: 'https://example.com/post2',
      isFavorite: true,
      tags: ['manga', 'character', 'art'],
      createdAt: undefined,
      updatedAt: undefined
    },
    {
      id: BigInt(3),
      sourceId: BigInt(2),
      deviceIds: [BigInt(2)],
      paths: ['/path/to/image3.jpg'],
      downloadUrl: 'https://example.com/image3.jpg',
      width: 1080,
      height: 1920,
      filesize: 612890,
      thumbnailPath: '/thumbnails/image3_thumb.jpg',
      imagePath: '/images/image3.jpg',
      postAuthor: 'Artist3',
      postAuthorUrl: 'https://example.com/artist3',
      postUrl: 'https://example.com/post3',
      isFavorite: false,
      tags: ['portrait', 'digital', 'fantasy'],
      createdAt: undefined,
      updatedAt: undefined
    }
  ]
  
  // Generate more dummy images with varied dimensions and dates
  const generateRandomDate = (daysBack: number) => {
    const now = new Date()
    const randomDays = Math.floor(Math.random() * daysBack)
    const date = new Date(now.getTime() - (randomDays * 24 * 60 * 60 * 1000))
    return {
      seconds: Math.floor(date.getTime() / 1000),
      nanos: 0
    }
  }

  const getRandomDimensions = () => {
    const aspectRatios = [
      { width: 1920, height: 1080 }, // 16:9
      { width: 1080, height: 1920 }, // 9:16 portrait
      { width: 1200, height: 1200 }, // 1:1 square
      { width: 1600, height: 900 },  // 16:9 wide
      { width: 800, height: 1200 },  // 2:3 portrait
      { width: 1400, height: 900 },  // ~16:10
      { width: 900, height: 1600 },  // tall portrait
    ]
    return aspectRatios[Math.floor(Math.random() * aspectRatios.length)]
  }

  const generateRandomTitle = () => {
    const titles = [
      'Mountain Sunset',
      'City Lights at Night',
      'Ocean Waves',
      'Forest Path',
      'Desert Landscape',
      'Autumn Leaves',
      'Winter Wonderland',
      'Spring Flowers',
      'Urban Architecture',
      'Wildlife Portrait',
      'Abstract Art',
      'Vintage Car',
      'Coffee Shop Morning',
      'Beach Sunset',
      'Northern Lights',
      'Rain Drops',
      'Star Field',
      'Garden View',
      'Street Art',
      'Peaceful Lake'
    ]
    return titles[Math.floor(Math.random() * titles.length)]
  }

  const images = Array.from({ length: 50 }, (_, index) => {
    const baseImage = dummyImages[index % dummyImages.length]
    const dimensions = getRandomDimensions()
    return {
      ...baseImage,
      id: BigInt(index + 1),
      imagePath: `/demo/image${index + 1}.jpg`,
      thumbnailPath: `/demo/thumb${index + 1}.jpg`,
      width: dimensions.width,
      height: dimensions.height,
      createdAt: generateRandomDate(30),
      title: generateRandomTitle()
    }
  })
  
  let selectedImage: Image | null = null

  function formatDate(timestamp?: any): string {
    if (!timestamp) return 'Unknown date'
    const date = new Date(timestamp.seconds * 1000)
    return date.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  }

  function groupImagesByDate(images: Image[]) {
    const groups: { [key: string]: Image[] } = {}
    
    images.forEach(image => {
      const dateKey = formatDate(image.createdAt)
      if (!groups[dateKey]) {
        groups[dateKey] = []
      }
      groups[dateKey].push(image)
    })
    
    return Object.entries(groups).sort(([a], [b]) => {
      const dateA = new Date(a).getTime()
      const dateB = new Date(b).getTime()
      return dateB - dateA // Sort newest first
    })
  }

  $: groupedImages = groupImagesByDate(images)
  
  function handleImageClick(event) {
    selectedImage = event.detail.image
  }
  
  function handleModalClose() {
    selectedImage = null
  }
</script>

<div class="gallery-container">
  <header class="gallery-header">
    <h1>Library</h1>
    <div class="controls">
      <div class="search-bar">
        <input type="text" placeholder="Search images..." />
      </div>
      <div class="view-options">
        <button class="view-btn active">Grid</button>
        <button class="view-btn">List</button>
      </div>
    </div>
  </header>
  
  <div class="gallery-content">
    {#each groupedImages as [date, dayImages]}
      <MasonryLayout 
        images={dayImages} 
        dateHeader={date}
        on:imageClick={handleImageClick}
      />
    {/each}
  </div>
</div>

<ImageModal 
  image={selectedImage}
  isOpen={!!selectedImage}
  on:close={handleModalClose}
/>

<style>
  .gallery-container {
    padding: 1rem;
    height: 100%;
    overflow: auto;
  }
  
  .gallery-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    padding: 0 0.5rem;
  }
  
  .gallery-header h1 {
    margin: 0;
    font-size: 2rem;
    font-weight: 600;
    color: #ffffff;
  }
  
  .controls {
    display: flex;
    gap: 1rem;
    align-items: center;
  }
  
  .search-bar input {
    padding: 0.5rem 1rem;
    background-color: #3d3d3d;
    border: 1px solid #5d5d5d;
    border-radius: 6px;
    color: #ffffff;
    width: 250px;
  }
  
  .search-bar input:focus {
    outline: none;
    border-color: #646cff;
  }
  
  .view-options {
    display: flex;
    gap: 0.5rem;
  }
  
  .view-btn {
    padding: 0.5rem 1rem;
    background-color: #3d3d3d;
    border: 1px solid #5d5d5d;
    border-radius: 6px;
    color: #aaa;
    cursor: pointer;
    transition: all 0.2s ease;
  }
  
  .view-btn:hover, .view-btn.active {
    background-color: #646cff;
    color: #ffffff;
    border-color: #646cff;
  }
  
  .gallery-content {
    padding: 0.5rem;
  }
</style>