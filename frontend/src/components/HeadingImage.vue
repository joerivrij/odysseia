<template>
    <div>
        <div class="topArea">
        <header>
            <h1>εξέταση</h1>
            <h2>Hoofdstuk {{this.$parent.selectedChapter}}</h2>
        </header>
        </div>
    </div>
</template>

<script>
    export default {
        name: "Quiz",
        data() {
        }
    }
</script>

<style scoped>

    ol {
        padding: 0;
        text-align: center;
        list-style-type: decimal;
        margin-left: 15em;
        width: 200px;
    }

    .topArea {
        padding: 5%;
        text-align: center;
        position: relative;
        opacity: 0.9;
        background-image: url("../assets/images/parthenon.jpg");
        background-size: cover;
        background-position: center center;
        background-attachment: fixed;
        width: 100%;
        height: 100%;
        z-index: -1;
    }

    /*!* On smaller screens, where height is less than 450px, change the style of the sidebar (less padding and a smaller font size) *!*/
    /*@media screen and (max-height: 450px) {*/
    /*    .sidenav {padding-top: 15px;}*/
    /*    .sidenav a {font-size: 18px;}*/
    /*}*/

</style>