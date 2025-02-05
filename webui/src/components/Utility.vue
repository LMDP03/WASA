<script>
import {RouterLink} from 'vue-router';

export default {
    props: {
        show: boolean,
        users: Array,
        title: String,
    },
    data() {
        return {

            error: "",

            searchName: "",

            searchnameValid: new RegExp('ˆ.*?$'),

            filteredUsers: [],
        };        
    },
    methods: {
        close() {
            this.searchName = "";
            window.location.reload();
            this.$emit('close');
        },
        async filterUsers() {
            this.error = "";

            this.filteredUsers = this.users;
            if (this.searchName.length < 3 || this.searchName.length > 16 || this.searchnameValid.test(this.searchName)) {
                this.error = "Invalid search name.";
                this.filteredUsers = [];
                return;
            }
            if (this.title == "search") {
                try {
                    const url = '/users/${sessionStorage.userId}/others?';

                }
            }
        }
    }

},
</script>