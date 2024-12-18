<!-- <template>
  <div>
    <form @submit.prevent="uploadImages">
      <input type="file" multiple @change="handleFileChange" />
      <button type="submit">Upload Images</button>
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';

const images = ref([]);
const handleFileChange = (event) => {
  images.value = event.target.files;
};

const UPLOAD_IMAGES = gql`
  mutation UploadImages($images: [String!]!) {
    uploadImages(images: $images) {
      image_urls
    }
  }
`;

const { mutate: uploadImagesMutation } = useMutation(UPLOAD_IMAGES);

const uploadImages = async () => {
  const imageStrings = Array.from(images.value).map(file => {
    // Convert file to Base64 or handle as needed
    return file.name; // Replace with actual file conversion logic
  });

  const { data } = await uploadImagesMutation({
    images: imageStrings,
  });

  console.log(data.uploadImages.image_urls);
};
</script> -->
