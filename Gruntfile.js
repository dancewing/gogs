
module.exports = function (grunt) {
    grunt.initConfig({
        pkg: grunt.file.readJSON('package.json'),

        sass: {
            options: {
                includePaths: [
                    'node_modules/foundation-sites/scss',
                    'node_modules/font-awesome/scss',
                    'node_modules/sass-flex-mixin/',
                    'node_modules/angularjs-datepicker/src/css/'
                ]
            },
            dist: {
                files: {
                    'public/css/app.css': 'src/scss/app.scss'
                }
            }
        },

		less : {

        	dev : {
				options: {
					paths: ['public/less']
				},
				files: {
					'public/css/gogs.css': 'public/less/gogs.less'
				}
			}

		},

        concat: {
            dist: {
                src: [
                    "src/**/*.module.js",
                    "src/**/**!(.module).js"
                ],
                dest: "public/js/board/app.js"
            }
        },

        uglify: {
            dist: {
                files: {
                    "public/js/board/app.min.js": ["public/js/board/app.js"]
                }
            }
        },

		clean : {

		},

        copy: {
            main: {
                files: [
                    {
                        flatten: false,
                        expand: true,
                        cwd: 'node_modules/twemoji/2/svg/',
                        src: ['**/*.svg'],
                        dest: 'public/img/twemoji/svg/',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/twemoji/2/twemoji.min.js'],
                        dest: 'public/js/board/twemoji.min.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/markdown-it-emoji/dist/markdown-it-emoji.min.js'],
                        dest: 'public/js/board/markdown-it-emoji.min.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/markdown-it/dist/markdown-it.js'],
                        dest: 'public/js/board/markdown-it.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-storage/dist/angular-storage.js'],
                        dest: 'public/js/board/angular-storage.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-lodash/angular-lodash.js'],
                        dest: 'public/js/board/angular-lodash.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/lodash/lodash.js'],
                        dest: 'public/js/board/lodash.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/topbar/topbar.js'],
                        dest: 'public/js/board/topbar.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/dropdownToggle/dropdownToggle.js'],
                        dest: 'public/js/board/dropdownToggle.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-file-upload/dist/angular-file-upload.js'],
                        dest: 'public/js/board/angular-file-upload.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/position/position.js'],
                        dest: 'public/js/board/position.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/typeahead/typeahead.js'],
                        dest: 'public/js/board/typeahead.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/bindHtml/bindHtml.js'],
                        dest: 'public/js/board/bindHtml.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/mediaQueries/mediaQueries.js'],
                        dest: 'public/js/board/mediaQueries.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/tabs/tabs.js'],
                        dest: 'public/js/board/tabs.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-loading-bar/build/loading-bar.js'],
                        dest: 'public/js/board/loading-bar.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-loading-bar/build/loading-bar.css'],
                        dest: 'public/css/loading-bar.css',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-mm-foundation/src/transition/transition.js'],
                        dest: 'public/js/board/transition.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/ng-sortable/dist/ng-sortable.js'],
                        dest: 'public/js/board/ng-sortable.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular/angular.js'],
                        dest: 'public/js/board/angular.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angular-ui-router/release/angular-ui-router.min.js'],
                        dest: 'public/js/board/angular-ui-router.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/ng-sortable/dist/ng-sortable.min.css'],
                        dest: 'public/css/ng-sortable.min.css',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        src: ['node_modules/angularjs-datepicker/dist/angular-datepicker.min.js'],
                        dest: 'public/js/board/angularjs-datepicker.min.js',
                        filter: 'isFile'
                    },
                    {
                        flatten: false,
                        expand: true,
                        cwd: 'node_modules/angular-mm-foundation/template',
                        src: '**',
                        dest: 'public/template/',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        expand: true,
                        cwd: 'node_modules/foundation-sites/js/foundation/',
                        src: '**',
                        dest: 'public/js/board',
                        filter: 'isFile'
                    },
                    {
                        flatten: true,
                        expand: true,
                        cwd: 'node_modules/font-awesome/fonts/',
                        src: '**',
                        dest: 'public/fonts/',
                        filter: 'isFile'
                    },
                    {
                        flatten: false,
                        expand: true,
                        cwd: 'src/',
                        src: ['**/*.js'],
                        dest: 'public/js/board/',
                        filter: 'isFile'
                    },
                    {
                        flatten: false,
                        expand: true,
                        cwd: 'src/',
                        src: ['**/*.html'],
                        dest: 'public/assets/html/',
                        filter: 'isFile'
                    }
        ]
    }
},

watch: {
    grunt: {
        files: ['Gruntfile.js'],
            tasks: ['sass', 'copy']
        },

        sass: {
            files: 'src/scss/*.scss',
            tasks: ['sass']
        },

		less: {
			files: 'public/less/*.less',
			tasks: ['less:dev']
		},

        copy: {
            files: ['src/**/*.js', 'src/**/*.html'],
            tasks: ['copy']
        },

        concat: {
            files: ['src/**/*.js'],
            tasks: ['concat']
        },

        uglify: {
            files: ['public/js/board/app.js'],
            tasks: ['uglify']
        }
    }
});

grunt.loadNpmTasks('grunt-sass');
grunt.loadNpmTasks('grunt-contrib-watch');
grunt.loadNpmTasks('grunt-contrib-copy');
grunt.loadNpmTasks('grunt-contrib-concat');
grunt.loadNpmTasks('grunt-contrib-uglify');

grunt.loadNpmTasks('grunt-contrib-less');

grunt.registerTask('build', ['sass', 'less:dev', 'copy', 'concat', 'uglify']);
grunt.registerTask('default', ['build', 'watch']);
};
