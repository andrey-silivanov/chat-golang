import {createRouter, createWebHistory} from 'vue-router';
import MainView from '../views/MainView.vue'
import LoginView from '../views/LoginView'
import ForbiddenView from "@/views/ForbiddenView";
import NotFoundView from "@/views/NotFoundView";
import RegisterView from "@/views/RegisterView";

const router = createRouter({
    hashbang: false,
    history: createWebHistory(),
    routes: [
        {
            path: '/',
            name: 'main',
            component: MainView,
            meta: {
                auth: true
            }
        },
        {
            path: '/register',
            name: 'register',
            component: RegisterView,
            meta: {
                auth: false
            }
        },
        {
            path: '/login',
            name: 'auth-login',
            component: LoginView,
            meta: {
                auth: false
            }
        },
        {
            path: '/403',
            name: '403',
            component: ForbiddenView,
        },
        {
            path: '/404',
            name: '404',
            component: NotFoundView,
        }
    ]
});

export default (app) => {
    app.router = router;

    app.use(router);
}